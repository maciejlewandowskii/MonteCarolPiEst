# Montecarlo Pi Estimator

A Go-based graphical application that estimates the value of Pi using the Monte Carlo method, built with the Fyne UI toolkit.

## Installation & Running

Download the latest release for your operating system from the [Releases](../../releases) page.

### 🍏 macOS Users: "App is damaged and can't be opened" Error

If you are on macOS and you see an error stating that the app **"is damaged and can't be opened. You should move it to the Trash"**, this is due to Apple's Gatekeeper security feature. 

**To fix this issue, run the following command in your terminal to remove the quarantine attribute:**

```bash
xattr -cr MontecarloPiEstimator.app
```

After running this command, you can double-click the app to launch it normally.

## Building from Source

To build the application yourself, ensure you have Go 1.24 or later installed:

```bash
git clone https://github.com/maciejlewandowskii/MonteCarolPiEst.git
cd MonteCarolPiEst
go mod download
go run main.go
```
