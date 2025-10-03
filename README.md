# luamerge

A CLI tool to merge Lua table files for Ragnarok Online, using JSON-based job configuration.

## ✨ Features

- ✅ **Deep recursive merge** of Lua tables with any nesting level
- ✅ **Multiple jobs** configurable in a single settings.json file
- ✅ **Multiple tables** per job - process several tables in a single file
- ✅ **Flexible paths** - support for relative and absolute paths
- ✅ **Embedded template** - no external dependencies
- ✅ **Robust validation** with descriptive error messages
- ✅ **Organized output** with full control over file destinations
- ✅ **Format agnostic** - handles explicit/implicit indices and different string key formats automatically

## 📦 Installation

### Quick Install (Linux/macOS)

Install the latest version:

```bash
curl -sSL https://raw.githubusercontent.com/zhaori96/luamerge/main/install.sh | bash
```

Install a specific version:

```bash
curl -sSL https://raw.githubusercontent.com/zhaori96/luamerge/main/install.sh | VERSION=v1.0.0 bash
```

Install to a custom directory:

```bash
curl -sSL https://raw.githubusercontent.com/zhaori96/luamerge/main/install.sh | INSTALL_DIR=/usr/local/bin bash
```

Or download and run the install script:

```bash
wget https://raw.githubusercontent.com/zhaori96/luamerge/main/install.sh
chmod +x install.sh
./install.sh

# Install specific version
VERSION=v1.0.0 ./install.sh

# Custom install directory
INSTALL_DIR=/custom/path ./install.sh
```

### Manual Installation

Download the latest binary from the [releases page](https://github.com/zhaori96/luamerge/releases) and add it to your PATH.

### Build from Source

```bash
git clone https://github.com/zhaori96/luamerge.git
cd luamerge
go build -o luamerge ./cmd/cli
```

## 🚀 Quick Start

1. Create an `input/` folder in the executable directory
2. Place your Lua files inside `input/`
3. Create a `settings.json` inside `input/`
4. Run:

```bash
luamerge                    # Uses input/ by default
luamerge --inputs myfolder/  # Uses custom folder
```

## 📝 Configuration (settings.json)

The `settings.json` file defines merge **jobs** and **global options**. Each job specifies:
- Which files to process (base and source)
- Where to save the result (output)
- Which tables and fields to merge (tables)
- Job-specific options (optional)

### Basic Structure

```json
{
  "options": {
    "keepUnmergedItems": false
  },
  "jobs": [
    {
      "name": "Job Name (optional)",
      "base": "base_file.lua",
      "source": "source_file.lua",
      "output": "output_file.lua",
      "options": {
        "keepUnmergedItems": true
      },
      "tables": {
        "TableName": {
          "field1": true,
          "field2": true
        }
      }
    }
  ]
}
```

### Options

#### `keepUnmergedItems` (boolean)

**Global (options)** or **per Job (job.options)**:

- `true`: **Preserves complete original file** (comments, variables, functions, other tables)
  - Only replaces tables specified in `tables`
  - Maintains original file order

- `false` (default): **Only specified tables**
  - Generates file only with tables configured in `tables`
  - Previous behavior

**Hierarchy**: Job options > Global options > Default (false)

### Complete Example

```json
{
  "options": {
    "keepUnmergedItems": false
  },
  "jobs": [
    {
      "name": "StateIcon PT-BR (Apenas tabela)",
      "base": "stateiconinfo.lua",
      "source": "stateiconinfo_latam.lua",
      "output": "stateiconinfo_ptbr.lua",
      "tables": {
        "StateIconList": {
          "descript": true
        }
      }
    },
    {
      "name": "Game Data (Arquivo completo preservado)",
      "base": "game_data.lua",
      "source": "game_data_ptbr.lua",
      "output": "../output/game_final.lua",
      "options": {
        "keepUnmergedItems": true
      },
      "tables": {
        "QuestInfoList": {
          "Title": true,
          "Description": true
        },
        "achievement_tbl": {
          "title": true,
          "content": {
            "summary": true,
            "details": true
          }
        }
      }
    }
  ]
}
```

**In this example:**
- Job 1 uses global (`keepUnmergedItems: false`) → Only `StateIconList` in output
- Job 2 overrides with (`keepUnmergedItems: true`) → Complete file with selective merge

### Merge Rules

#### 1. Complete Field Merge
```json
"field": true
```
Replaces the entire field with the value from the source file.

#### 2. Selective Nested Merge
```json
"field": {
  "subfield1": true,
  "subfield2": true
}
```
Inside `field`, only `subfield1` and `subfield2` are replaced.

#### 3. Deep Merge (Multiple Levels)
```json
"field": {
  "level1": {
    "level2": {
      "level3": true
    }
  }
}
```
Supports any depth of nesting.

#### 4. Complete Table Replacement
```json
"TableName": true
```
Replaces the entire table with the source version.

### File Paths

#### Input Files (base and source)
- **Always relative to input/ folder**
- `"base": "file.lua"` → `input/file.lua`
- `"base": "subfolder/file.lua"` → `input/subfolder/file.lua`

#### Output Files
- **Filename only**: goes to `output/` next to `input/`
  ```json
  "output": "result.lua"  // → output/result.lua
  ```
- **Relative path**: can use `../`
  ```json
  "output": "../custom/result.lua"  // → custom/result.lua
  ```
- **Absolute path**: uses exact path
  ```json
  "output": "/tmp/result.lua"  // → /tmp/result.lua
  ```

## 📂 Project Structure

```
luamerge/
├── cmd/cli/              # CLI application
│   └── main.go
├── internal/
│   ├── parser/          # Lua file parser (AST)
│   │   ├── parser.go
│   │   └── table.go
│   ├── config/          # Configuration loading and validation
│   │   └── settings.go
│   ├── merger/          # Recursive merge logic
│   │   ├── merger.go
│   │   └── result.go
│   ├── preservation/    # Text-based preservation
│   │   └── textmerge.go
│   └── template/        # Embedded Lua template
│       ├── template.go
│       └── lua.gotmpl
├── input/               # 📥 User working folder
│   ├── settings.json    # Job configuration
│   └── *.lua            # Lua files to process
├── output/              # 📤 Generated results (automatic)
│   └── *.lua
├── install.sh           # Installation script
└── README.md
```

## ⚙️ How It Works

1. **Loading**: Reads `settings.json` from input/ folder
2. **Parser**: Analyzes Lua files using AST (Abstract Syntax Tree)
3. **Merge**: Applies merge rules recursively for each job
4. **Generation**: Uses embedded template to create output files
5. **Output**: Saves results as configured in each job

## 💡 Complete Usage Example

```bash
# 1. Initial structure
mkdir input
cd input

# 2. Copy Lua files
cp /path/stateiconinfo.lua .
cp /path/stateiconinfo_latam.lua .

# 3. Create settings.json
cat > settings.json << 'EOF'
{
  "jobs": [
    {
      "name": "StateIcon Translation",
      "base": "stateiconinfo.lua",
      "source": "stateiconinfo_latam.lua",
      "output": "stateiconinfo_final.lua",
      "tables": {
        "StateIconList": {
          "descript": true
        }
      }
    }
  ]
}
EOF

# 4. Run
cd ..
luamerge

# 5. Result in output/stateiconinfo_final.lua
```

## 🎯 Use Cases

### 1. Single File Translation
```json
{
  "jobs": [
    {
      "name": "StateIcon Translation",
      "base": "stateiconinfo.lua",
      "source": "stateiconinfo_translated.lua",
      "output": "stateiconinfo_final.lua",
      "tables": {
        "StateIconList": { "descript": true }
      }
    }
  ]
}
```

### 2. Multiple Files in Batch
```json
{
  "jobs": [
    {
      "name": "StateIcon",
      "base": "stateicon.lua",
      "source": "stateicon_translated.lua",
      "output": "stateicon_final.lua",
      "tables": { "StateIconList": { "descript": true } }
    },
    {
      "name": "Quest",
      "base": "quest.lua",
      "source": "quest_translated.lua",
      "output": "quest_final.lua",
      "tables": { "QuestInfoList": { "Title": true, "Description": true } }
    }
  ]
}
```

### 3. Multiple Tables in Same File
```json
{
  "jobs": [
    {
      "name": "Complete Game Data",
      "base": "gamedata.lua",
      "source": "gamedata_translated.lua",
      "output": "gamedata_final.lua",
      "tables": {
        "StateIconList": { "descript": true },
        "QuestInfoList": { "Title": true, "Description": true },
        "achievement_tbl": { "title": true, "content": { "summary": true } }
      }
    }
  ]
}
```

## 🔧 Development

```bash
# Build
go build -o luamerge ./cmd/cli

# Run tests
go test ./...

# Production build
go build -ldflags="-s -w" -o luamerge ./cmd/cli
```

## 📄 License

MIT
