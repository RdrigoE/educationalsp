local client = vim.lsp.start_client {
    name = "educationalsp",
    cmd = { "/home/reusebio/projets/educationalsp/main" }
}

if not client then
    vim.notify "hey you didnt the client thing good"
    return
end

vim.api.nvim_create_autocmd("FileType", {
    pattern = "markdown",
    callback = function()
        vim.lsp.buf_attach_client(0, client)
    end,
})
