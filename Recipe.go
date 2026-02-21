//go:build gobake
package bake_recipe

import (
	"fmt"
	"runtime"
	"github.com/fezcode/gobake"
)

func Run(bake *gobake.Engine) error {
	if err := bake.LoadRecipeInfo("recipe.piml"); err != nil {
		return err
	}

	bake.Task("build", "Builds the binary for the current platform", func(ctx *gobake.Context) error {
		ctx.Log("Building %s v%s for %s/%s...", bake.Info.Name, bake.Info.Version, runtime.GOOS, runtime.GOARCH)

		osName := runtime.GOOS
		archName := runtime.GOARCH

		err := ctx.Mkdir("build")
		if err != nil {
			return err
		}

		ldflags := fmt.Sprintf("-X main.Version=%s", bake.Info.Version)

		output := "build/" + bake.Info.Name + "-" + osName + "-" + archName
		if osName == "windows" {
			output += ".exe"
		}

		ctx.Env = []string{
			"CGO_ENABLED=0", // Pure Go audio
			"GOOS=" + osName,
			"GOARCH=" + archName,
		}
		
		err = ctx.Run("go", "build", "-ldflags", ldflags, "-o", output, ".")
		return err
	})

	bake.Task("clean", "Removes build artifacts", func(ctx *gobake.Context) error {
		return ctx.Remove("build")
	})

	return nil
}
