# atlas.radio 📻🏜️

![Banner Image](./banner-image.png)

**atlas.radio** is a high-fidelity, world radio receiver for the Atlas Suite. Tune into thousands of live stations across the globe directly from your terminal. Built with a pure Go audio engine, it requires zero external dependencies like FFmpeg or VLC.

Part of the **Atlas Suite**, it features a rugged industrial TUI interface with real-time signal telemetry.

![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)

## ✨ Features

- 📡 **Global Reach:** Access the open Radio Browser API to search for any station worldwide.
- ⚙️ **Pure Go Audio:** Uses `oto/v3` and `go-mp3` for native, cross-platform audio playback.
- 📟 **Atlas Industrial TUI:** High-fidelity Amber terminal aesthetic with a minimalist footprint.
- 📊 **Signal Telemetry:** Real-time frequency data and station metadata display.
- 🔍 **Frequency Search:** Instant fuzzy search for stations by name, tag, or country.
- 📜 **Signal Catalog:** Smooth scrolling interface to browse through 50+ discovered signals.

## 🚀 Installation

### Recommended: Via Atlas Hub
The easiest way to install is using the central hub:
```bash
atlas.hub
```
Select `atlas.radio` from the list and confirm.

### From Source
```bash
git clone https://github.com/fezcode/atlas.radio
cd atlas.radio
gobake build
```

## ⌨️ Usage

Simply run the binary to start receiving signals:
```bash
./atlas.radio
```

### 🕹️ Controls
| Key | Action |
|-----|--------|
| `↑/↓` / `j/k` | **Navigate:** Move through the signal list. |
| `Enter` | **Tune In:** Connect to the selected radio station. |
| `s` / `/` | **Search:** Open the frequency search input. |
| `p` | **Power Off:** Stop the current audio stream. |
| `q` / `Ctrl+C` | **Exit:** Shut down the receiver. |

---

## 🛠️ Technical Stack

- **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling:** [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **Audio Driver:** [Oto v3](https://github.com/ebitengine/oto)
- **MP3 Decoder:** [go-mp3](https://github.com/hajimehoshi/go-mp3)
- **API:** [Radio Browser](https://www.radio-browser.info/)

## 📄 License
MIT License - see [LICENSE](LICENSE) for details.
