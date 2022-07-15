export default (_: Request): Response => {
    return new Response(`Zenith is running!`, {status: 200})
}