package domain

import "time"

// LotteryResult representa o resultado de um sorteio da loteria
type LotteryResult struct {
	ID            string    `json:"id"`
	GameType      string    `json:"gameType"`
	DrawNumber    int       `json:"drawNumber"`
	DrawDate      string    `json:"drawDate"`
	Numbers       []int     `json:"numbers"`
	Winners       int       `json:"winners"`
	Prize         float64   `json:"prize"`
	ProcessedAt   time.Time `json:"processedAt"`
	Source        string    `json:"source"`
}

// DownstreamPayload é o payload enviado para o microserviço downstream
type DownstreamPayload struct {
	Status   string          `json:"status"`
	Results  []LotteryResult `json:"results"`
	ErrorMsg string          `json:"errorMsg,omitempty"`
	Metadata Metadata        `json:"metadata"`
}

// Metadata contém informações sobre a requisição
type Metadata struct {
	ProcessedAt   time.Time `json:"processedAt"`
	SourceService string    `json:"sourceService"`
	TotalRecords  int       `json:"totalRecords"`
	RequestID     string    `json:"requestId"`
}

// WebhookPayload representa um webhook externo
type WebhookPayload struct {
	GameType   string    `json:"gameType"`
	DrawNumber int       `json:"drawNumber"`
	DrawDate   string    `json:"drawDate"`
	Numbers    []int     `json:"numbers"`
	Timestamp  time.Time `json:"timestamp"`
}

// CaixaAPIResponse representa a resposta da API da Caixa
type CaixaAPIResponse struct {
	Acumulado                          bool     `json:"acumulado"`
	DataApuracao                       string   `json:"dataApuracao"`
	DataProximoConcurso                string   `json:"dataProximoConcurso"`
	DezenasSorteadasOrdemSorteio       []string `json:"dezenasSorteadasOrdemSorteio"`
	ExibirDetalhamentoPorCidade        bool     `json:"exibirDetalhamentoPorCidade"`
	ID                                 *int     `json:"id"`
	IndicadorConcursoEspecial          int      `json:"indicadorConcursoEspecial"`
	ListaDezenas                       []string `json:"listaDezenas"`
	ListaDezenasSegundoSorteio         []string `json:"listaDezenasSegundoSorteio"`
	ListaMunicipioUFGanhadores         []string `json:"listaMunicipioUFGanhadores"`
	ListaRateioPremio                  []struct {
		DescricaoFaixa      string  `json:"descricaoFaixa"`
		Faixa               int     `json:"faixa"`
		NumeroDeGanhadores  int     `json:"numeroDeGanhadores"`
		ValorPremio         float64 `json:"valorPremio"`
	} `json:"listaRateioPremio"`
	ListaResultadoEquipeEsportiva      interface{} `json:"listaResultadoEquipeEsportiva"`
	LocalSorteio                       string      `json:"localSorteio"`
	NomeMunicipioUFSorteio             string      `json:"nomeMunicipioUFSorteio"`
	NomeTimeCoracaoMesSorte            string      `json:"nomeTimeCoracaoMesSorte"`
	Numero                             int         `json:"numero"`
	NumeroConcursoAnterior             int         `json:"numeroConcursoAnterior"`
	NumeroConcursoFinal05              int         `json:"numeroConcursoFinal_0_5"`
	NumeroConcursoProximo              int         `json:"numeroConcursoProximo"`
	NumeroJogo                         int         `json:"numeroJogo"`
	Observacao                         string      `json:"observacao"`
	PremiacaoContingencia              interface{} `json:"premiacaoContingencia"`
	TipoJogo                           string      `json:"tipoJogo"`
	TipoPublicacao                     int         `json:"tipoPublicacao"`
	UltimoConcurso                     bool        `json:"ultimoConcurso"`
	ValorArrecadado                    float64     `json:"valorArrecadado"`
	ValorAcumuladoConcurso05           float64     `json:"valorAcumuladoConcurso_0_5"`
	ValorAcumuladoConcursoEspecial     float64     `json:"valorAcumuladoConcursoEspecial"`
	ValorAcumuladoProximoConcurso      float64     `json:"valorAcumuladoProximoConcurso"`
	ValorEstimadoProximoConcurso       float64     `json:"valorEstimadoProximoConcurso"`
	ValorSaldoReservaGarantidora       float64     `json:"valorSaldoReservaGarantidora"`
	ValorTotalPremioFaixaUm            float64     `json:"valorTotalPremioFaixaUm"`
}

// ServiceInfo contém informações do serviço
type ServiceInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Uptime      string `json:"uptime"`
	Status      string `json:"status"`
}

// HealthResponse é a resposta do health check
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// MetricsResponse contém métricas do serviço
type MetricsResponse struct {
	RequestsTotal   int64         `json:"requestsTotal"`
	RequestsSuccess int64         `json:"requestsSuccess"`
	RequestsError   int64         `json:"requestsError"`
	AverageLatency  time.Duration `json:"averageLatency"`
	CacheHits       int64         `json:"cacheHits"`
	CacheMisses     int64         `json:"cacheMisses"`
}

// ErrorResponse é a resposta padrão de erro
type ErrorResponse struct {
	Error      string `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Timestamp  string `json:"timestamp"`
	RequestID  string `json:"requestId"`
}

// PaginatedResponse é uma resposta paginada
type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Total int         `json:"total"`
}
