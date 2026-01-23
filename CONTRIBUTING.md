# Contributing

Obrigado por considerar contribuir para o Lottery Caixa Service!

## Como Contribuir

### 1. Relatando Bugs

Se encontrou um bug, por favor abra uma issue descrevendo:
- Como reproduzir o problema
- Comportamento esperado
- Comportamento observado
- Versão do Go e OS

### 2. Sugerindo Melhorias

Para sugerir melhorias:
1. Use o título descritivo
2. Forneça uma descrição clara da solução
3. Inclua exemplos práticos

### 3. Pull Requests

#### Setup de Desenvolvimento

```bash
# 1. Fork o repositório
git clone https://github.com/seu-username/lottery-caixa-service.git
cd lottery-caixa-service

# 2. Criar branch para feature
git checkout -b feature/sua-feature

# 3. Instalar dependências
go mod download

# 4. Fazer alterações
# ... edite seus arquivos ...

# 5. Executar testes
make test

# 6. Formatar código
make fmt

# 7. Commit
git add .
git commit -m "Add: descrição da feature"

# 8. Push
git push origin feature/sua-feature

# 9. Abrir Pull Request
```

#### Padrões de Código

- Siga o style guide do Go
- Use `goimports` para formatação
- Escreva testes para novas features
- Mantenha cobertura de testes > 80%
- Adicione documentação para funções públicas

#### Mensagens de Commit

```
<type>: <subject>

<body>

<footer>
```

**Types:**
- `feat`: Nova feature
- `fix`: Bug fix
- `docs`: Documentação
- `style`: Formatação
- `refactor`: Refatoração
- `perf`: Performance
- `test`: Testes
- `chore`: Build, deps, etc

**Exemplos:**
```
feat: add retry logic to Caixa API client
fix: resolve cache expiration issue
docs: update deployment guide for Cloud Run
```

### 4. Processo de Review

1. Mínimo 1 aprovação necessária
2. Testes devem passar
3. Cobertura não pode diminuir
4. Documentação deve ser atualizada

### 5. Guia de Estilo

#### Go

```go
// Sempre use nomes descritivos
func FetchLotteryResults(ctx context.Context, gameType string) (*domain.LotteryResult, error) {
    // Implementação
}

// Documente funções públicas
// FetchLotteryResults busca resultados da loteria da Caixa com retry automático.
// Retorna erro se falhar após max retries.

// Use error wrapping
if err != nil {
    return nil, fmt.Errorf("erro ao buscar resultados: %w", err)
}

// Prefira explicitismo
if err != nil {
    // Não: return err
    // Sim:
    return fmt.Errorf("ao processar payload: %w", err)
}
```

#### Testes

```go
func TestFetchLotteryResults(t *testing.T) {
    // Arrange
    cfg := &config.Config{...}
    svc := service.NewLotteryService(cfg)
    
    // Act
    result, err := svc.FetchLotteryResults(context.Background(), "lotofacil")
    
    // Assert
    if err != nil {
        t.Fatalf("esperava sucesso, obteve erro: %v", err)
    }
    
    if result == nil {
        t.Error("resultado não deve ser nil")
    }
}
```

#### Documentação

```go
// Package service contém a lógica de negócio principal.
package service

// LotteryService gerencia operações de loteria.
type LotteryService struct {
    // campos privados
    cfg *config.Config
}

// NewLotteryService cria uma nova instância de LotteryService.
// Retorna erro se a configuração for inválida.
func NewLotteryService(cfg *config.Config) *LotteryService {
    // implementação
}
```

### 6. Estrutura de Diretórios

Mantenha a estrutura consistente:

```
lottery-caixa-service/
├── cmd/              # Aplicação principal
├── internal/         # Pacotes privados
│   ├── domain/       # Modelos
│   ├── service/      # Lógica de negócio
│   ├── cache/        # Cache
│   ├── ratelimit/    # Rate limiting
│   └── http/         # HTTP handlers
├── config/           # Configuração
├── scripts/          # Scripts úteis
├── examples/         # Exemplos
└── docs/            # Documentação
```

### 7. Checklist para PR

- [ ] Testes escritos e passando
- [ ] Cobertura mantida ou melhorada
- [ ] Código formatado (`make fmt`)
- [ ] Linter passa (`make lint`)
- [ ] README atualizado se necessário
- [ ] CHANGELOG atualizado
- [ ] Mensagem de commit clara
- [ ] Commits amassados (squash) se necessário

### 8. Workflow de Release

1. Atualizar versão em `config.go`
2. Atualizar CHANGELOG.md
3. Criar tag: `git tag v1.0.0`
4. Push tag: `git push origin v1.0.0`
5. GitHub Actions cria release automaticamente

### 9. Principais Áreas

#### Métricas e Monitoramento
- Adicionar novas métricas prometheus
- Melhorar health checks
- Adicionar distributed tracing

#### Performance
- Otimizações de cache
- Melhorar rate limiting
- Reduzir latência

#### Funcionalidades
- Suporte para mais loterias
- Integração com outras APIs
- Novos tipos de webh ooks

### 10. Perguntas?

Abra uma issue ou entre em contato!

---

Obrigado por contribuir! 🎉
