package config

import "tool"

func init() {
   tool.Add("app", tool.StrMap{
      "name":  tool.Env("app.APP_NAME", "shop"),
      "env":   tool.Env("app.APP_ENV", "production"),
      "debug": tool.Env("app.APP_DEBUG",false),
      "port":  tool.Env("app.APP_PORT","8089"),
   })
}

