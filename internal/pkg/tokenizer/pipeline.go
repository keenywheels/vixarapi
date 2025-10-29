package tokenizer

// Pipeline represents a sequence of processing stages for tokens.
type Pipeline struct {
	initStage PipelineStage
}

// Run processes the input tokens through the pipeline and returns the processed tokens.
func (p *Pipeline) Run(tokens []Token) []Token {
	if p.initStage == nil {
		return tokens
	}

	return p.initStage.Execute(tokens)
}

// PipelineBuilder helps in constructing a Pipeline by adding stages.
type PipelineBuilder struct {
	stages []PipelineStage
}

// NewPipelineBuilder creates a new instance of PipelineBuilder.
func NewPipelineBuilder() *PipelineBuilder {
	return &PipelineBuilder{}
}

// AddStages adds multiple stages to the pipeline.
func (b *PipelineBuilder) AddStages(stages ...PipelineStage) *PipelineBuilder {
	b.stages = append(b.stages, stages...)

	return b
}

// Build constructs the Pipeline with the added stages.
func (b *PipelineBuilder) Build() *Pipeline {
	initStage := connectPipelineStages(b.stages...)

	return &Pipeline{
		initStage: initStage,
	}
}

// connectPipelineStages links the provided stages in sequence and returns the first stage.
func connectPipelineStages(stages ...PipelineStage) PipelineStage {
	var (
		firstStage PipelineStage
		lastStage  PipelineStage
	)

	for _, stage := range stages {
		if firstStage == nil {
			firstStage = stage
		}

		if lastStage != nil {
			lastStage.SetNext(stage)
		}

		lastStage = stage
	}

	return firstStage
}
