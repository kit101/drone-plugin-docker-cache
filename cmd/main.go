package main

import (
	"github.com/joho/godotenv"
	plugin "github.com/kit101/drone-plugins-docker-cache"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

var (
	version = "unknown"
	usage   = `docker data root缓存插件，复制已有docker data root目录到registry path并处理.dockerignore`
)

func main() {
	// Load env-file if it exists first
	if env := os.Getenv("PLUGIN_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := cli.NewApp()
	app.Name = "docker cache plugin"
	app.Usage = usage
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "storage-path",
			Usage:    "docker data root 目录",
			EnvVar:   "PLUGIN_STORAGE_PATH",
			Required: true,
		},
		cli.StringFlag{
			Name:   "src",
			Usage:  "docker data root缓存目录.\n该值存在会复制缓存目录到storage-path. \n若storage-path是在workingDir下，则还会在的${workingDir}/.dockerignore中追加storage-path的相对路径.",
			EnvVar: "PLUGIN_SRC",
		},
		cli.StringSliceFlag{
			Name:   "dockerignores",
			Usage:  ".dockerignore中额外写入的忽略路径",
			EnvVar: "PLUGIN_DOCKERIGNORES",
		},

		/* TODO 未实现
		cli.StringFlag{
			Name:   "dest",
			Usage:  "docker data root 保存目录",
			EnvVar: "PLUGIN_DEST",
		},
		cli.StringFlag{
			Name:   "dest-type",
			Usage:  "dest type: dir or tar",
			EnvVar: "PLUGIN_DEST_TYPE",
		},
		*/
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	p := plugin.Plugin{
		StoragePath:   c.String("storage-path"),
		Src:           c.String("src"),
		Dockerignores: c.StringSlice("dockerignores"),

		/*
			TODO 未实现
			Dest:          c.String("desc"),
			DestType:      c.String("dest-type"),
		*/

	}
	return p.Exec()
}
