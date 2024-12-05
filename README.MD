
# goroutines-advanced-logging

Este repositório demonstra o uso avançado de **goroutines**, controle de concorrência e boas práticas em logs no desenvolvimento backend com Go. 

**Aviso:** Este código é parte de um projeto privado e proprietário. Foi publicado com a permissão explícita do proprietário original para fins de exibição como exemplo em meu portfólio público.

A funcionalidade apresentada faz parte de um projeto maior e foi adaptada para servir como exemplo prático de habilidades específicas, incluindo:

- **Concorrência e Paralelismo**:
  - Utilização de **goroutines** para execução paralela de tarefas.
  - Controle de concorrência com **sync.WaitGroup** e **canais**.
- **Integração com MongoDB**:
  - Paginação e filtragem de dados diretamente do banco de dados.
  - Uso de **filtros dinâmicos** para consultas otimizadas.
- **Tratamento de Logs**:
  - Registro detalhado de erros e eventos relevantes para diagnóstico.
  - Mensagens estruturadas para facilitar a análise e depuração.
- **Tratamento de Erros Robusto**:
  - Captura de `panics` para evitar interrupções inesperadas.
  - Mensagens de erro amigáveis para usuários e desenvolvedores.
- **Formatação da Resposta**:
  - Combinação de dados de múltiplas fontes para compor uma resposta JSON rica e informativa.

---

## 🚀 Funcionalidade Principal

A função apresentada no repositório realiza a listagem de mudanças de plantão em um sistema fictício. Ela inclui:

1. Consulta paginada ao MongoDB.
2. Filtro dinâmico para seleção de registros.
3. Processamento paralelo para busca de nomes associados aos IDs retornados.
4. Montagem de uma resposta JSON detalhada e ordenada.

---

## 💡 Destaques Técnicos

- **Controle de Concorrência:** Uso de `sync.WaitGroup` para gerenciar múltiplas goroutines.
- **Filas e Buffers:** Utilização de canais para comunicação entre goroutines.
- **Estrutura Modular:** Código limpo e separado em responsabilidades específicas.
- **Escalabilidade:** Capacidade de processar grandes volumes de dados de maneira eficiente.

---

## 🛠️ Pré-requisitos

- **Go**: Versão 1.18 ou superior.
- **MongoDB**: Configuração local ou em nuvem.
- **Bibliotecas Go**:
  - [Fiber](https://gofiber.io/): Framework web rápido e minimalista.
  - [MongoDB Driver](https://github.com/mongodb/mongo-go-driver): Para integração com MongoDB.

---

## 📦 Como Executar

1. Clone o repositório:
   ```bash
   git clone https://github.com/seu-usuario/goroutines-advanced-logging.git
   cd goroutines-advanced-logging
   ```

2. Instale as dependências:
   ```bash
   go mod tidy
   ```

3. Configure o MongoDB no arquivo de configuração.

4. Execute o servidor:
   ```bash
   go run main.go
   ```

5. Faça uma requisição para a rota:
   ```bash
   curl http://localhost:3000/change-duties?page=1&limit=10
   ```

---

## 📝 Licença

Este repositório está disponível apenas para fins educacionais e como exemplo de portfólio. A publicação deste trecho foi autorizada pelo proprietário original.

---

Criado por [Seu Nome](https://github.com/seu-usuario).
