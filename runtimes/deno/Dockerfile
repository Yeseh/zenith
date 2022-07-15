FROM denoland/deno:latest as base

ARG ENTRYPOINT=zenith-main.ts
ARG PORT=80

EXPOSE ${PORT} 

ENV ZENITH_APP_ENTRYPOINT=$ENTRYPOINT
ENV ZENITH_APP_PORT=$PORT

WORKDIR /app
COPY . ./
RUN rm ./functions/ping.ts
RUN deno run --allow-read --allow-write gen-functionfile.ts
RUN rm gen-functionfile.ts
RUN deno cache $ENTRYPOINT

CMD deno run --allow-net --allow-env $ZENITH_APP_ENTRYPOINT