FROM ronmi/mingo
ADD forgejo-mcp /forgejo-mcp
ENTRYPOINT ["/forgejo-mcp"]
CMD ["stdio"]
