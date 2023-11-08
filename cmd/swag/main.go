package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/nguyennm96/swaggo-v3"
	"github.com/nguyennm96/swaggo-v3/format"
	"github.com/nguyennm96/swaggo-v3/gen"
)

const (
	searchDirFlag         = "dir"
	excludeFlag           = "exclude"
	generalInfoFlag       = "generalInfo"
	propertyStrategyFlag  = "propertyStrategy"
	outputFlag            = "output"
	outputTypesFlag       = "outputTypes"
	parseVendorFlag       = "parseVendor"
	parseDependencyFlag   = "parseDependency"
	markdownFilesFlag     = "markdownFiles"
	codeExampleFilesFlag  = "codeExampleFiles"
	parseInternalFlag     = "parseInternal"
	generatedTimeFlag     = "generatedTime"
	requiredByDefaultFlag = "requiredByDefault"
	parseDepthFlag        = "parseDepth"
	instanceNameFlag      = "instanceName"
	overridesFileFlag     = "overridesFile"
	parseGoListFlag       = "parseGoList"
	quietFlag             = "quiet"
	tagsFlag              = "tags"
	parseExtensionFlag    = "parseExtension"
	templateDelimsFlag    = "templateDelims"
	openAPIVersionFlag    = "v3.1"
	packageName           = "packageName"
	collectionFormatFlag  = "collectionFormat"
)

var initFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:    quietFlag,
		Aliases: []string{"q"},
		Usage:   "Make the logger quiet.",
	},
	&cli.StringFlag{
		Name:    generalInfoFlag,
		Aliases: []string{"g"},
		Value:   "main.go",
		Usage:   "Go file path in which 'swagger general API Info' is written",
	},
	&cli.StringFlag{
		Name:    searchDirFlag,
		Aliases: []string{"d"},
		Value:   "./",
		Usage:   "Directories you want to parse,comma separated and general-info file must be in the first one",
	},
	&cli.StringFlag{
		Name:  excludeFlag,
		Usage: "Exclude directories and files when searching, comma separated",
	},
	&cli.StringFlag{
		Name:    propertyStrategyFlag,
		Aliases: []string{"p"},
		Value:   swaggo.CamelCase,
		Usage:   "Property Naming Strategy like " + swaggo.SnakeCase + "," + swaggo.CamelCase + "," + swaggo.PascalCase,
	},
	&cli.StringFlag{
		Name:    outputFlag,
		Aliases: []string{"o"},
		Value:   "./docs",
		Usage:   "Output directory for all the generated files(swagger.json, swagger.yaml and docs.go)",
	},
	&cli.StringFlag{
		Name:    outputTypesFlag,
		Aliases: []string{"ot"},
		Value:   "go,json,yaml",
		Usage:   "Output types of generated files (docs.go, swagger.json, swagger.yaml) like go,json,yaml",
	},
	&cli.BoolFlag{
		Name:  parseVendorFlag,
		Usage: "Parse go files in 'vendor' folder, disabled by default",
	},
	&cli.BoolFlag{
		Name:    parseDependencyFlag,
		Aliases: []string{"pd"},
		Usage:   "Parse go files inside dependency folder, disabled by default",
	},
	&cli.StringFlag{
		Name:    markdownFilesFlag,
		Aliases: []string{"md"},
		Value:   "",
		Usage:   "Parse folder containing markdown files to use as description, disabled by default",
	},
	&cli.StringFlag{
		Name:    codeExampleFilesFlag,
		Aliases: []string{"cef"},
		Value:   "",
		Usage:   "Parse folder containing code example files to use for the x-codeSamples extension, disabled by default",
	},
	&cli.BoolFlag{
		Name:  parseInternalFlag,
		Usage: "Parse go files in internal packages, disabled by default",
	},
	&cli.BoolFlag{
		Name:  generatedTimeFlag,
		Usage: "Generate timestamp at the top of docs.go, disabled by default",
	},
	&cli.IntFlag{
		Name:  parseDepthFlag,
		Value: 100,
		Usage: "Dependency parse depth",
	},
	&cli.BoolFlag{
		Name:  requiredByDefaultFlag,
		Usage: "Set validation required for all fields by default",
	},
	&cli.StringFlag{
		Name:  instanceNameFlag,
		Value: "",
		Usage: "This parameter can be used to name different swagger document instances. It is optional.",
	},
	&cli.StringFlag{
		Name:  overridesFileFlag,
		Value: gen.DefaultOverridesFile,
		Usage: "File to read global type overrides from.",
	},
	&cli.BoolFlag{
		Name:  parseGoListFlag,
		Value: true,
		Usage: "Parse dependency via 'go list'",
	},
	&cli.StringFlag{
		Name:  parseExtensionFlag,
		Value: "",
		Usage: "Parse only those operations that match given extension",
	},
	&cli.StringFlag{
		Name:    tagsFlag,
		Aliases: []string{"t"},
		Value:   "",
		Usage:   "A comma-separated list of tags to filter the APIs for which the documentation is generated.Special case if the tag is prefixed with the '!' character then the APIs with that tag will be excluded",
	},
	&cli.BoolFlag{
		Name:  openAPIVersionFlag,
		Value: false,
		Usage: "Generate OpenAPI V3.1 spec",
	},
	&cli.StringFlag{
		Name:    templateDelimsFlag,
		Aliases: []string{"td"},
		Value:   "",
		Usage:   "Provide custom delimeters for Go template generation. The format is leftDelim,rightDelim. For example: \"[[,]]\"",
	},
	&cli.StringFlag{
		Name:  packageName,
		Value: "",
		Usage: "A package name of docs.go, using output directory name by default (check `--output` option)",
	},
	&cli.StringFlag{
		Name:    collectionFormatFlag,
		Aliases: []string{"cf"},
		Value:   "csv",
		Usage:   "Set default collection format",
	},
}

func initAction(ctx *cli.Context) error {
	strategy := ctx.String(propertyStrategyFlag)

	switch strategy {
	case swaggo.CamelCase, swaggo.SnakeCase, swaggo.PascalCase:
	default:
		return fmt.Errorf("not supported %s propertyStrategy", strategy)
	}

	leftDelim, rightDelim := "{{", "}}"

	if ctx.IsSet(templateDelimsFlag) {
		delims := strings.Split(ctx.String(templateDelimsFlag), ",")
		if len(delims) != 2 {
			return fmt.Errorf("exactly two template delimeters must be provided, comma separated")
		} else if delims[0] == delims[1] {
			return fmt.Errorf("template delimiters must be different")
		}
		leftDelim, rightDelim = strings.TrimSpace(delims[0]), strings.TrimSpace(delims[1])
	}

	outputTypes := strings.Split(ctx.String(outputTypesFlag), ",")
	if len(outputTypes) == 0 {
		return fmt.Errorf("no output types specified")
	}
	logger := log.New(os.Stdout, "", log.LstdFlags)
	if ctx.Bool(quietFlag) {
		logger = log.New(io.Discard, "", log.LstdFlags)
	}

	collectionFormat := swaggo.TransToValidCollectionFormat(ctx.String(collectionFormatFlag))
	if collectionFormat == "" {
		return fmt.Errorf("not supported %s collectionFormat", ctx.String(collectionFormat))
	}

	return gen.New().Build(&gen.Config{
		SearchDir:           ctx.String(searchDirFlag),
		Excludes:            ctx.String(excludeFlag),
		ParseExtension:      ctx.String(parseExtensionFlag),
		MainAPIFile:         ctx.String(generalInfoFlag),
		PropNamingStrategy:  strategy,
		OutputDir:           ctx.String(outputFlag),
		OutputTypes:         outputTypes,
		ParseVendor:         ctx.Bool(parseVendorFlag),
		ParseDependency:     ctx.Bool(parseDependencyFlag),
		MarkdownFilesDir:    ctx.String(markdownFilesFlag),
		ParseInternal:       ctx.Bool(parseInternalFlag),
		GeneratedTime:       ctx.Bool(generatedTimeFlag),
		RequiredByDefault:   ctx.Bool(requiredByDefaultFlag),
		CodeExampleFilesDir: ctx.String(codeExampleFilesFlag),
		ParseDepth:          ctx.Int(parseDepthFlag),
		InstanceName:        ctx.String(instanceNameFlag),
		OverridesFile:       ctx.String(overridesFileFlag),
		ParseGoList:         ctx.Bool(parseGoListFlag),
		Tags:                ctx.String(tagsFlag),
		LeftTemplateDelim:   leftDelim,
		RightTemplateDelim:  rightDelim,
		PackageName:         ctx.String(packageName),
		Debugger:            logger,
		GenerateOpenAPI3Doc: ctx.Bool(openAPIVersionFlag),
		CollectionFormat:    collectionFormat,
	})
}

func main() {
	fmt.Println("Swag version: ", swaggo.Version)
	app := cli.NewApp()
	app.Version = swaggo.Version
	app.Usage = "Automatically generate RESTful API documentation with Swagger 3.0 for Go."
	app.Commands = []*cli.Command{
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "Create docs.go",
			Action:  initAction,
			Flags:   initFlags,
		},
		{
			Name:    "fmt",
			Aliases: []string{"f"},
			Usage:   "format swag comments",
			Action: func(c *cli.Context) error {
				searchDir := c.String(searchDirFlag)
				excludeDir := c.String(excludeFlag)
				mainFile := c.String(generalInfoFlag)

				return format.New().Build(&format.Config{
					SearchDir: searchDir,
					Excludes:  excludeDir,
					MainFile:  mainFile,
				})
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    searchDirFlag,
					Aliases: []string{"d"},
					Value:   "./",
					Usage:   "Directories you want to parse,comma separated and general-info file must be in the first one",
				},
				&cli.StringFlag{
					Name:  excludeFlag,
					Usage: "Exclude directories and files when searching, comma separated",
				},
				&cli.StringFlag{
					Name:    generalInfoFlag,
					Aliases: []string{"g"},
					Value:   "main.go",
					Usage:   "Go file path in which 'swagger general API Info' is written",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}