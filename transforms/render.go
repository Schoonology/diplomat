package transforms

import (
	"regexp"

	"github.com/testdouble/diplomat/parsers"
	"github.com/testdouble/diplomat/scripting"
)

var templateChunk *regexp.Regexp

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

func renderBodies(test *parsers.Test) error {
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

// RenderTemplates renders all the Headers and Bodies in all the Tests in the
// given `spec`.
func RenderTemplates(spec *parsers.Spec) error {
	for _, test := range spec.Tests {
		if err := renderAllHeaders(test.Request.Headers); err != nil {
			return err
		}
		if err := renderAllHeaders(test.Response.Headers); err != nil {
			return err
		}
		if err := renderBodies(&test); err != nil {
			return err
		}
	}

	return nil
}

// RenderTemplatesImmutable renders all the Headers and Bodies in all the Tests in the
// given `spec`.
// TODO: Consolidate all of this
func RenderTemplatesImmutable(spec *parsers.Spec) (*parsers.Spec, error) {
	for _, test := range spec.Tests {
		if err := renderAllHeaders(test.Request.Headers); err != nil {
			return nil, err
		}
		if err := renderAllHeaders(test.Response.Headers); err != nil {
			return nil, err
		}
		if err := renderBodies(&test); err != nil {
			return nil, err
		}
	}

	return spec, nil
}

// RenderTemplatesStream updates a stream of specs using RenderTemplates.
// given `specChannel`.
func RenderTemplatesStream(specChannel chan *parsers.Spec, errorChannel chan error) (chan *parsers.Spec, chan error) {
	c := make(chan *parsers.Spec)
	e := make(chan error)

	go func() {
		var spec *parsers.Spec
		select {
		case spec = <-specChannel:
		case err := <-errorChannel:
			e <- err
			return
		}

		updatedSpec, err := RenderTemplatesImmutable(spec)
		if err != nil {
			e <- err
		}

		c <- updatedSpec
	}()

	return c, e
}
