package template

var ZenithTomlConfigTemplate = `[app]
name = "{{.AppName}}"
runtime = "{{.Runtime}}"

[[function]]
name = "ping"
path = "./functions/ping.ts"
route = "/api/functions"
`

var DenoPingFuncTemplate = `export default (_: Request): Response => {
    return new Response( "Zenith is running!", {status: 200} )
};
`

var PingFuncRuntimeMap = map[string]string{
	"deno": DenoPingFuncTemplate,
}
