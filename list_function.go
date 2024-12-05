// list_function.go
// Este arquivo contém uma implementação demonstrativa da função `list`, projetada para
// destacar o uso de recursos avançados da linguagem Go. A função serve como um exemplo
// didático, abordando conceitos importantes como:
// 
// - **Concorrência e Goroutines**: Execução paralela de tarefas para melhorar a eficiência.
// - **Canais e WaitGroups**: Sincronização e comunicação entre goroutines.
// - **Tratamento de Erros Robusto**: Captura e manuseio de erros com mensagens informativas.
// - **Interação com MongoDB**: Consultas paginadas e filtradas em uma base de dados.
// - **Logs Estruturados**: Registro detalhado de eventos e erros para depuração e monitoramento.
// 
// Este código foi extraído de um projeto maior para exibição como exemplo público e
// busca demonstrar boas práticas e o poder do Go em cenários reais.

func (c *ChangeDutyController) List(ctx *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recuperado de um pânico:", r)
			err := ctx.Status(500).JSON(fiber.Map{
				"message": "Ocorreu um erro interno no servidor",
			})
			if err != nil {
				return
			}
		}
	}()

	if c.MongoClient == nil {
		log.Println("ChangeDutyController: MongoClient is nil")
		return ctx.Status(500).JSON(fiber.Map{
			"error": "MongoClient is nil",
		})
	}

	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit

	report := ctx.Query("report", "false") == "true"

	filter := bson.M{}
	if report {
		now := time.Now()
		startDate := now.AddDate(0, 0, -30)
		filter["created_at"] = bson.M{"$gte": startDate}
	}

	collection := c.MongoClient.GetCollection(config.MaisCuidadoDB, config.ChangeDutyCollection)

	var changeDuties []models.ChangeDuty
	optionsPersonal := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{"start_time", -1}})
	cursor, err := collection.Find(ctx.Context(), filter, optionsPersonal)
	if err != nil {
		log.Println("Error on finding all change duties:", err)
		return ctx.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := cursor.All(ctx.Context(), &changeDuties); err != nil {
		log.Println("Error on cursor.All:", err)
		return ctx.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	clientIDs := make([]string, 0)
	oldCaregiverIDs := make([]string, 0)
	newCaregiverIDs := make([]string, 0)
	userIDs := make([]string, 0)

	for _, changeDuty := range changeDuties {
		if changeDuty.ClientID != "" {
			clientIDs = append(clientIDs, changeDuty.ClientID)
		}

		if changeDuty.OldCaregiverID != "" {
			if _, err := primitive.ObjectIDFromHex(changeDuty.OldCaregiverID); err == nil {
				oldCaregiverIDs = append(oldCaregiverIDs, changeDuty.OldCaregiverID)
			} else {
				log.Println("OldCaregiverID inválido ou mal formado:", changeDuty.OldCaregiverID)
				oldCaregiverIDs = append(oldCaregiverIDs, changeDuty.OldCaregiverID)
			}
		}

		if changeDuty.NewCaregiverID != "" {
			if _, err := primitive.ObjectIDFromHex(changeDuty.NewCaregiverID); err == nil {
				newCaregiverIDs = append(newCaregiverIDs, changeDuty.NewCaregiverID)
			} else {
				log.Println("NewCaregiverID inválido ou mal formado:", changeDuty.NewCaregiverID)
				newCaregiverIDs = append(newCaregiverIDs, changeDuty.NewCaregiverID)
			}
		}

		if changeDuty.UserID != "" {
			if _, err := primitive.ObjectIDFromHex(changeDuty.UserID); err == nil {
				userIDs = append(userIDs, changeDuty.UserID)
			} else {
				log.Println("UserID inválido ou mal formado:", changeDuty.UserID)
				userIDs = append(userIDs, changeDuty.UserID)
			}
		}
	}

	clientNamesChan := make(chan map[string]string)
	oldCaregiverNamesChan := make(chan map[string]string)
	newCaregiverNamesChan := make(chan map[string]string)
	userNamesChan := make(chan map[string]string)
	errorChan := make(chan error)
	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		defer wg.Done()
		clientNames, err := c.ClientController.GetClientNamesByIDs(clientIDs)
		if err != nil {
			log.Println("Erro ao buscar nomes dos clientes:", err)
			clientNames = make(map[string]string)
			for _, id := range clientIDs {
				clientNames[id] = "Falha ao obter Nome: " + id
			}
		}
		clientNamesChan <- clientNames
	}()

	go func() {
		defer wg.Done()
		oldCaregiverNames, err := c.CaregiverController.GetCaregiverNamesByIDs(oldCaregiverIDs)
		if err != nil {
			log.Println("Erro ao buscar nomes dos cuidadores antigos:", err)
			oldCaregiverNames = make(map[string]string)
			for _, id := range oldCaregiverIDs {
				if _, err := primitive.ObjectIDFromHex(id); err != nil {
					oldCaregiverNames[id] = "Falha ao obter Nome: ID inválido - " + id
				} else {
					oldCaregiverNames[id] = "Falha ao obter Nome: " + id
				}
			}
		}
		oldCaregiverNamesChan <- oldCaregiverNames
	}()

	go func() {
		defer wg.Done()
		newCaregiverNames, err := c.CaregiverController.GetCaregiverNamesByIDs(newCaregiverIDs)
		if err != nil {
			log.Println("Erro ao buscar nomes dos novos cuidadores:", err)
			newCaregiverNames = make(map[string]string)
			for _, id := range newCaregiverIDs {
				if _, err := primitive.ObjectIDFromHex(id); err != nil {
					newCaregiverNames[id] = "Falha ao obter Nome: ID inválido - " + id
				} else {
					newCaregiverNames[id] = "Falha ao obter Nome: " + id
				}
			}
		}
		newCaregiverNamesChan <- newCaregiverNames
	}()

	go func() {
		defer wg.Done()
		userNames, err := c.UserController.GetUserNamesByIDs(userIDs)
		if err != nil {
			log.Println("Erro ao buscar nomes dos usuários:", err)
			userNames = make(map[string]string)
			for _, id := range userIDs {
				if _, err := primitive.ObjectIDFromHex(id); err != nil {
					userNames[id] = "Falha ao obter Nome: ID inválido - " + id
				} else {
					userNames[id] = "Falha ao obter Nome: " + id
				}
			}
		}
		userNamesChan <- userNames
	}()

	go func() {
		wg.Wait()
		close(clientNamesChan)
		close(oldCaregiverNamesChan)
		close(newCaregiverNamesChan)
		close(userNamesChan)
		close(errorChan)
	}()

	var clientNames, oldCaregiverNames, newCaregiverNames, userNames map[string]string
	for {
		select {
		case names := <-clientNamesChan:
			clientNames = names
		case names := <-oldCaregiverNamesChan:
			oldCaregiverNames = names
		case names := <-newCaregiverNamesChan:
			newCaregiverNames = names
		case names := <-userNamesChan:
			userNames = names
		case err := <-errorChan:
			if err != nil {
				log.Println("Error on getting names:", err)
				return ctx.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		}
		if clientNames != nil && oldCaregiverNames != nil && newCaregiverNames != nil && userNames != nil {
			break
		}
	}

	var changeDutyResponses []models.ChangeDutyResponse
	for _, changeDuty := range changeDuties {
		changeDutyResponse := models.ChangeDutyResponse{
			ID:               changeDuty.ID,
			ClientID:         changeDuty.ClientID,
			ClientName:       clientNames[changeDuty.ClientID],
			OldCaregiverID:   changeDuty.OldCaregiverID,
			OldCaregiverName: oldCaregiverNames[changeDuty.OldCaregiverID],
			NewCaregiverID:   changeDuty.NewCaregiverID,
			NewCaregiverName: newCaregiverNames[changeDuty.NewCaregiverID],
			UserID:           changeDuty.UserID,
			UserName:         userNames[changeDuty.UserID],
			StartTime:        changeDuty.StartTime,
			EndTime:          changeDuty.EndTime,
			DutySituation:    changeDuty.DutySituation,
			DutyDescription:  changeDuty.DutyDescription,
			CreatedAt:        changeDuty.CreatedAt,
			UpdatedAt:        changeDuty.UpdatedAt,
			UUID:             changeDuty.UUID,
		}

		changeDutyResponses = append(changeDutyResponses, changeDutyResponse)
	}

	sort.Slice(changeDutyResponses, func(i, j int) bool {
		return changeDutyResponses[j].StartTime.Before(changeDutyResponses[i].StartTime)
	})

	return ctx.JSON(changeDutyResponses)
}
