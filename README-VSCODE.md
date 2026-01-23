# Guia VS Code - Lottery Caixa Service

Guia completo para desenvolver e debugar o projeto no Visual Studio Code.

## 🎯 Configuração Inicial

### 1. Abrir o Projeto

```bash
cd /Users/victorhxavier/Documents/lottery-caixa-service
code .
```

### 2. Instalar Extensões Recomendadas

Ao abrir o projeto, o VS Code irá sugerir extensões. Clique em **"Install All"** ou instale manualmente:

**Essenciais:**
- **Go** (`golang.go`) - Suporte completo para Go
- **Docker** (`ms-azuretools.vscode-docker`) - Gerenciar containers
- **REST Client** (`humao.rest-client`) - Testar API no VS Code

**Recomendadas:**
- **GitLens** (`eamodio.gitlens`) - Git superpowers
- **Error Lens** (`usernamehw.errorlens`) - Erros inline
- **DotENV** (`mikestead.dotenv`) - Syntax highlight para .env
- **Prettier** (`esbenp.prettier-vscode`) - Formatação

### 3. Configurar Go Tools

Pressione `Cmd+Shift+P` (Mac) ou `Ctrl+Shift+P` (Windows/Linux) e digite:

```
Go: Install/Update Tools
```

Selecione **TODAS** as ferramentas e clique em **OK**.

## 🚀 Executando o Projeto

### Método 1: Debug Panel (F5)

1. Vá para o painel **Run and Debug** (ícone de play com bug na barra lateral)
2. Selecione uma das configurações:
   - **Launch Service** - Apenas o serviço principal
   - **Launch Processor Only** - Apenas o processador
   - **Service + Processor (Debug Both)** - Ambos simultaneamente
3. Pressione **F5** ou clique no botão play verde

### Método 2: Terminal Integrado

Pressione `` Ctrl+` `` para abrir o terminal e execute:

```bash
# Executar serviço
go run ./cmd/main.go

# Ou usar Make
make run-dev
```

### Método 3: Tasks (Cmd+Shift+B)

Pressione `Cmd+Shift+B` (Mac) ou `Ctrl+Shift+B` (Windows/Linux) para ver as tasks:

- **build** - Compilar o projeto
- **run** - Executar o serviço
- **test** - Rodar testes
- **test-coverage** - Testes com cobertura
- **docker: compose up** - Iniciar com Docker Compose

## 🐛 Debugando

### Breakpoints

1. Clique na margem esquerda do editor (antes do número da linha) para adicionar breakpoint
2. Pressione **F5** para iniciar debug
3. Use os controles:
   - **F5** - Continue
   - **F10** - Step Over
   - **F11** - Step Into
   - **Shift+F11** - Step Out
   - **Shift+F5** - Stop

### Configurações de Debug Disponíveis

#### 1. Launch Service
Debug do serviço principal com breakpoints

#### 2. Launch Service + Processor
Debug de ambos os serviços simultaneamente

#### 3. Launch Processor Only
Debug apenas do processador downstream

#### 4. Test Current File
Debug dos testes do arquivo atual

#### 5. Test Current Package
Debug de todos os testes do pacote

#### 6. Test All
Debug de todos os testes do projeto

### Debug Variables

Durante o debug, você pode:
- Ver **VARIABLES** no painel esquerdo
- Ver **WATCH** para expressões customizadas
- Ver **CALL STACK** para rastreamento
- Usar **DEBUG CONSOLE** para avaliar expressões Go

## 🧪 Testando a API

### Método 1: REST Client (Recomendado)

1. Abrir arquivo `.vscode/api-test.http`
2. Certificar que o serviço está rodando
3. Clicar em **"Send Request"** acima de cada requisição
4. Ver resposta no painel lateral

**Atalhos:**
- `Cmd+Alt+R` (Mac) ou `Ctrl+Alt+R` (Windows/Linux) - Enviar requisição

### Método 2: Terminal

```bash
# Executar script de testes
./test-api.sh
```

### Método 3: cURL no Terminal Integrado

```bash
# Health check
curl http://localhost:8080/health

# Buscar loteria
curl "http://localhost:8080/?gameType=lotofacil"
```

## 📝 Tasks Disponíveis

Pressione `Cmd+Shift+P` → `Tasks: Run Task`:

### Build & Run
- **build** - Compilar binário
- **run** - Executar serviço
- **start-processor** - Iniciar processador downstream
- **start-local** - Script interativo de inicialização

### Testing
- **test** - Rodar todos os testes
- **test-coverage** - Testes com relatório de cobertura
- **test-race** - Testes com race detector
- **test-api** - Testar endpoints da API

### Code Quality
- **format** - Formatar código com gofmt
- **lint** - Executar linter (golangci-lint)
- **go: tidy** - Limpar dependências

### Docker
- **docker: build** - Build da imagem Docker
- **docker: compose up** - Iniciar todos os serviços
- **docker: compose down** - Parar serviços
- **docker: compose logs** - Ver logs

### Utilities
- **clean** - Limpar artefatos de build
- **go: download dependencies** - Baixar dependências

## ⌨️ Atalhos Úteis

### Gerais
- `Cmd+P` - Quick Open (abrir arquivo)
- `Cmd+Shift+P` - Command Palette
- `Cmd+B` - Toggle Sidebar
- `` Cmd+` `` - Toggle Terminal
- `Cmd+Shift+E` - Explorer
- `Cmd+Shift+F` - Search

### Go Específicos
- `F12` - Go to Definition
- `Alt+F12` - Peek Definition
- `Shift+F12` - Find All References
- `F2` - Rename Symbol
- `Cmd+.` - Quick Fix
- `Cmd+Shift+O` - Go to Symbol

### Debug
- `F5` - Start/Continue Debug
- `F9` - Toggle Breakpoint
- `F10` - Step Over
- `F11` - Step Into
- `Shift+F5` - Stop Debug

### Testing
- `Cmd+Shift+B` - Run Build Task
- `Cmd+Shift+T` - Run Test Task (custom)

## 📂 Estrutura de Arquivos

```
.vscode/
├── launch.json          # Configurações de debug
├── tasks.json           # Tasks automatizadas
├── settings.json        # Configurações do workspace
├── extensions.json      # Extensões recomendadas
└── api-test.http        # Testes de API REST Client
```

## 🔧 Configurações Importantes

### Format on Save

O código será formatado automaticamente ao salvar (`.vscode/settings.json:26-29`):
- Go usa `goimports`
- JSON/YAML usa `prettier`
- Auto-organização de imports

### Go Tools

Configurado em `.vscode/settings.json:2-21`:
- `gopls` como Language Server
- `golangci-lint` para linting
- Race detection habilitado nos testes
- Code lens para testes e referências

### Auto Save

Arquivos são salvos automaticamente após 1 segundo de inatividade.

## 🎨 Temas Recomendados

Para melhor experiência:
- **Dark+** (padrão)
- **One Dark Pro**
- **Dracula Official**
- **Night Owl**

## 💡 Dicas de Produtividade

### 1. Multi-Cursor

`Alt+Click` para adicionar múltiplos cursores

### 2. Pesquisa Rápida

`Cmd+P` → Digite `@` para ver símbolos do arquivo atual

### 3. Terminal Split

`Cmd+\` no terminal para dividir

### 4. Zen Mode

`Cmd+K Z` para modo sem distrações

### 5. Source Control

`Cmd+Shift+G` para ver Git changes

### 6. Snippets Go

- `pkgm` - Package main
- `func` - Function
- `forr` - For range loop
- `if` - If statement
- `switch` - Switch statement

## 🐳 Docker no VS Code

### Gerenciar Containers

1. Clique no ícone Docker na sidebar
2. Veja containers, imagens, volumes, networks
3. Botão direito para ações (start, stop, logs, etc.)

### Docker Compose

- **Start**: Clique direito no `docker-compose.yml` → "Compose Up"
- **Stop**: Clique direito → "Compose Down"
- **Logs**: Clique direito no container → "View Logs"

## 🧪 Testes no VS Code

### Executar Testes

**Opção 1:** Code Lens
- Clique em "run test" acima de cada função de teste

**Opção 2:** Testing Panel
- Ícone de laboratório na sidebar
- Ver todos os testes
- Run/Debug individual ou em grupo

**Opção 3:** Terminal
```bash
go test -v ./...
```

### Ver Cobertura

Execute a task `test-coverage` e abra `coverage.html` no browser.

## 🔍 Troubleshooting

### Go Tools não funcionam

```bash
# Reinstalar ferramentas
Cmd+Shift+P → Go: Install/Update Tools
```

### Linter muito lento

Desabilite temporariamente em `.vscode/settings.json`:
```json
"go.lintOnSave": "off"
```

### Breakpoints não funcionam

1. Certifique-se que está em modo Debug (F5)
2. Recompile: `go clean -cache`
3. Reinicie VS Code

### Imports não organizando

```bash
# Instalar goimports
go install golang.org/x/tools/cmd/goimports@latest
```

### Port já em uso

```bash
# Matar processo na porta 8080
lsof -ti:8080 | xargs kill -9
```

## 📊 Status Bar

Informações úteis na barra inferior:
- **Go version** - Versão do Go
- **Go Test** - Status dos testes
- **GOPROXY** - Proxy do Go
- **Line/Column** - Posição no arquivo

## 🎯 Workflows Comuns

### Desenvolvimento Normal

1. Abrir VS Code
2. Pressione `F5` (Launch Service)
3. Edite o código
4. Salve (auto-format)
5. Veja mudanças refletidas

### Adicionar Nova Feature

1. Crie/edite arquivos
2. Adicione testes
3. Execute `test` task
4. Debug se necessário com breakpoints
5. Teste API com `.vscode/api-test.http`

### Debug de Problema

1. Adicione breakpoints
2. Pressione `F5`
3. Reproduza o problema
4. Inspecione variables
5. Use debug console para testar hipóteses

## 🔗 Links Úteis

- [VS Code Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.go)
- [VS Code Docs](https://code.visualstudio.com/docs)
- [Debugging in VS Code](https://code.visualstudio.com/docs/editor/debugging)
- [Go in VS Code](https://code.visualstudio.com/docs/languages/go)

## 💬 Dicas Finais

1. **Use Command Palette** (`Cmd+Shift+P`) - É seu melhor amigo
2. **Aprenda atalhos** - Economize tempo
3. **Personalize** - Ajuste `settings.json` ao seu gosto
4. **Use REST Client** - Muito mais rápido que Postman
5. **Debug sempre** - Breakpoints > Print statements
6. **Git integrado** - Não precisa sair do VS Code

---

**Happy Coding! 🚀**

Para dúvidas sobre o projeto, veja [README.md](README.md) ou [README-LOCAL.md](README-LOCAL.md)
