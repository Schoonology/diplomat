package transforms

import (
	"regexp"

	"github.com/testdouble/diplomat/builders"
	"github.com/testdouble/diplomat/scripting"
)

var templateChunk *regexp.Regexp

// TemplateRenderer renders all the Headers and Bodies in all the Tests in the
// provided channel.
type TemplateRenderer struct{}

func init() {
	templateChunk = regexp.MustCompilePOSIX("{{ ([^}]+) }}")
}

func renderTemplateBytes(src []byte) (dst []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	dst = templateChunk.ReplaceAllFunc(src, func(match []byte) []byte {
		submatch := templateChunk.FindSubmatch(match)
		pipelineSrc := submatch[1]

		result, err := scripting.RunPipeline(string(pipelineSrc))
		if err != nil {
			panic(err)
		}

		return []byte(result)
	})

	return
}

func renderTemplateString(src string) (dst string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	dst = templateChunk.ReplaceAllStringFunc(src, func(match string) string {
		submatch := templateChunk.FindStringSubmatch(match)
		pipelineSrc := submatch[1]

		result, err := scripting.RunPipeline(pipelineSrc)
		if err != nil {
			panic(err)
		}

		return result
	})

	return
}

func renderAllHeaders(headers map[string]string) error {
	for key, value := range headers {
		newValue, err := renderTemplateString(value)
		if err != nil {
			return err
		}

		headers[key] = newValue
	}

	return nil
}

func renderBodies(test *builders.Test) error {
	newBody, err := renderTemplateBytes(test.Request.Body)
	if err != nil {
		return err
	}

	test.Request.Body = newBody

	newBody, err = renderTemplateBytes(test.Response.Body)
	if err != nil {
		return err
	}

	test.Response.Body = newBody

	return nil
}

// Transform renders all the Headers and Bodies in a single Test.
func (*TemplateRenderer) Transform(test builders.Test) (builders.Test, error) {
	if err := renderAllHeaders(test.Request.Headers); err != nil {
		return test, err
	}
	if err := renderAllHeaders(test.Response.Headers); err != nil {
		return test, err
	}
	if err := renderBodies(&test); err != nil {
		return test, err
	}

	return test, nil
}

// TransformAll renders all the Headers and Bodies in all the Tests in the
// provided channel.
func (renderer *TemplateRenderer) TransformAll(tests chan builders.Test, errors chan error) chan builders.Test {
	output := make(chan builders.Test)

	go func() {
		for test := range tests {
			rendered, err := renderer.Transform(test)
			if err != nil {
				errors <- err
				return
			}

			output <- rendered
		}

		close(output)
	}()

	return output
}
