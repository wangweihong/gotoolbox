package waitgroup

type Output struct {
	Data interface{} `json:"data"`
	Err  error       `json:"err"`
}

type BatchOutput struct {
	Total   int       `json:"total"`
	Success int       `json:"success"`
	Fail    int       `json:"fail"`
	Results []*Output `json:"results"`
}

func (b *BatchOutput) SetOutput(data interface{}, err error) {
	b.Total += 1
	if err != nil {
		b.Fail += 1
		b.Results = append(b.Results, &Output{Data: data, Err: err})
		return
	}

	b.Success += 1
	b.Results = append(b.Results, &Output{Data: data, Err: err})
}

func (b *BatchOutput) Merge(c *BatchOutput) {
	if c == nil {
		return
	}

	b.Total += c.Total
	b.Success += c.Success
	b.Fail += c.Fail

	if c.Results != nil {
		b.Results = append(b.Results, c.Results...)
	}
}

func SetOutput(data interface{}, err error) *Output {
	var MyOutput Output
	MyOutput.Err = err
	if data == nil {
		MyOutput.Data = ""
	} else {
		MyOutput.Data = data
	}
	return &MyOutput
}
