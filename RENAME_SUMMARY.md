# Rename Summary: goluatbl → luamerge

## ✅ Renomeação Completa Realizada

O projeto foi completamente renomeado de `goluatbl` para `luamerge`.

### 📝 Arquivos Modificados:

#### **Go Module & Imports**
- ✅ `go.mod` - Module name: `goluatbl` → `luamerge`
- ✅ `cmd/cli/main.go` - Todos os imports atualizados
- ✅ `internal/merger/result.go` - Import atualizado
- ✅ `internal/merger/merger.go` - Import atualizado
- ✅ `internal/preservation/textmerge.go` - Import atualizado

#### **Binary & Commands**
- ✅ `cmd/cli/main.go` - Command name: `goluatbl` → `luamerge`
- ✅ Cobra command `Use` field atualizado
- ✅ Descrição e mensagens de output atualizadas

#### **Build & Release Configuration**
- ✅ `install.sh` - Todas as referências atualizadas:
  - REPO: `your-username/goluatbl` → `your-username/luamerge`
  - BINARY_NAME: `goluatbl` → `luamerge`
  - Mensagens de instalação atualizadas

- ✅ `.goreleaser.yml` - Configuração completa:
  - Build ID: `goluatbl` → `luamerge`
  - Binary name: `goluatbl` → `luamerge`
  - Archive ID: `goluatbl` → `luamerge`
  - Release repo: `goluatbl` → `luamerge`
  - Upload URL atualizada

#### **Documentation**
- ✅ `README.md` - Todas as 64 ocorrências substituídas
- ✅ `RELEASE_SETUP.md` - Todas as referências atualizadas
- ✅ `INSTALL_EXAMPLES.md` - Todas as referências atualizadas

#### **GitHub Actions**
- ✅ `.github/workflows/release.yml` - Sem mudanças necessárias (agnóstico ao nome)

### 🧪 Testes Realizados:

```bash
✅ go build -o luamerge ./cmd/cli
✅ ./luamerge --help
✅ ./luamerge --inputs input
```

**Resultado:** Todos os testes passaram com sucesso!

### 📊 Estatísticas da Renomeação:

- **Arquivos Go modificados:** 5
- **Arquivos de configuração modificados:** 2 (install.sh, .goreleaser.yml)
- **Arquivos de documentação modificados:** 3 (README.md, RELEASE_SETUP.md, INSTALL_EXAMPLES.md)
- **Total de arquivos alterados:** 10
- **Binário compilado:** 6.1 MB

### 🚀 Próximos Passos:

1. **Renomear o diretório do projeto (opcional):**
   ```bash
   cd ..
   mv goluatbl luamerge
   cd luamerge
   ```

2. **Atualizar repositório Git:**
   ```bash
   git add .
   git commit -m "refactor: rename project from goluatbl to luamerge"
   ```

3. **Criar novo repositório no GitHub (se ainda não existir):**
   - Criar repositório `luamerge` no GitHub
   - Atualizar remote: `git remote set-url origin https://github.com/YOUR_USERNAME/luamerge.git`

4. **Substituir YOUR_USERNAME nos arquivos:**
   ```bash
   sed -i 's/YOUR_USERNAME/seususuario/g' install.sh .goreleaser.yml README.md INSTALL_EXAMPLES.md RELEASE_SETUP.md
   ```

5. **Push e criar primeira release:**
   ```bash
   git push origin main
   git tag v1.0.0
   git push origin v1.0.0
   ```

### ✨ O Que Mudou:

**Antes:**
```bash
go build -o goluatbl ./cmd/cli
./goluatbl --inputs input
```

**Depois:**
```bash
go build -o luamerge ./cmd/cli
./luamerge --inputs input
```

**Instalação Antes:**
```bash
curl -sSL https://raw.githubusercontent.com/USER/goluatbl/main/install.sh | bash
```

**Instalação Depois:**
```bash
curl -sSL https://raw.githubusercontent.com/USER/luamerge/main/install.sh | bash
```

### 🎯 Checklist Final:

- [x] Module name atualizado
- [x] Imports internos atualizados
- [x] Binary name atualizado
- [x] Install script atualizado
- [x] GoReleaser config atualizado
- [x] README atualizado
- [x] Toda documentação atualizada
- [x] Build testado e funcionando
- [x] Execução testada e funcionando

## ✅ Renomeação Completa e Testada!

O projeto está 100% renomeado e funcional com o novo nome `luamerge`.
