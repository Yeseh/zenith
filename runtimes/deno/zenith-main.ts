import { serve } from "https://deno.land/std@0.148.0/http/server.ts";
import { functionMap } from "./functions.ts"

const PORT = Number(Deno.env.get("ZENITH_APP_PORT") ?? "80");
const BASE_URL = `http://localhost:${PORT}`;

for (const [key, _] of functionMap) {
  const msg = `Function: ${BASE_URL}/${key.pathname}`
  console.log(msg) 
}

const findHandler = (url: string) => {
  for (const [pattern, func] of functionMap) {
    if (pattern.exec(url)) {
      return func;
    }
  }
}

const handler = async (request: Request): Promise<Response> => {
  const handlerFunc = findHandler(request.url)
  if (!handlerFunc) {
    return new Response("Unknown function", { status: 404 })
  }

  return await handlerFunc(request)
};

await serve(handler, { port: PORT });