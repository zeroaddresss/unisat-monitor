# ⚡ Unisat Monitor

## Table of Contents

- [🚀 Introduction](#introduction)
- [🔍 Project Overview](#project-overview)
- [✨ Features](#features)
- [🏁 Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [🔑 API Key Information](#api-key-information)
  - [⚙️ Configuration](#configuration)
  - [📦 Installation and Usage](#installation-and-usage)
- [🎥 Demo](#demo)
- [🤝🏻 Contributing](#contributing)
- [⚠️ Disclaimer](#disclaimer)

## 🚀 Introduction

This project is an enhanced, **Go-implemented** version of the [unisat-monitor](link) repository.

This tool is designed to track the status of **[BRC20](https://www.brc20.guide/) collections** on [Unisat](https://unisat.io/market), by monitoring the floor price of collections and detecting price changes. It can offer **competitive advantages** for trading activities, by sending Discord webhook notifications when the bot detects a listing priced significantly below the current floor.

I rewrote the NodeJS codebase from scratch in Go to gain proficiency in **Golang programming**. Indeed, this project has been a great practice to understand Go's coding philosophy, core language constructs, and idiomatic Go practices.

## 🔍 Project Overview

Developed in April, during a period of high interest in trading BRC20 collections on Bitcoin, this tool has proven valuable in detecting and capitalizing on **pricing errors** (listings mistakenly priced too low) ahead of other traders. Initially, the ultimate goal of this project was to evolve it by:

1. Extending support to the highly anticipated [**Runes protocol**](https://docs.ordinals.com/runes.html)
2. Creating a **sniper bot** that would automatically purchase mispriced listings, without the need for any user interaction to manually purchase them

However, the interest in BRC20 tokens quickly faded, and I decided to stop because the product did not meet any market needs. Nonetheless, **pull requests** for implementing the mentioned functionalities are welcome and greatly appreciated.

## ✨ Features

This Go-based version offers enhanced functionality compared to its [NodeJS counterpart](link):

- **Parallel monitoring** of multiple collections leveraging goroutines
- Improved **logging**
- Enhanced **error handling**
- Cleaner and **well-structured** codebase

## 🏁 Getting Started

### Prerequisites

- **Golang** installed on your system
- **Unisat API key(s)**

### 🔑 API Key Information

The program relies on the **Unisat API** to retrieve data. While it's technically possible to make the bot work without API keys (using web scraping), this project uses official Unisat API keys for ethical and legal reasons.

Indeed, one or more API keys are required for proper functionality. Unisat's free tier allows up to **10,000 requests per day**. If multiple API keys are provided, the bot will shuffle and use them alternately to distribute requests evenly across the available keys.
To obtain an API key, you can contact Unisat via email or Telegram.

### ⚙️ Configuration

Before running the bot, user must populate the `config.json` file with the desired monitoring parameters, according to their preference.

### 📦 Installation and Usage

1. Verify that Go is correctly installed on your machine:

   ```sh
   go version
   ```

2. Clone the repository:

   ```sh
   git clone [repository-link]
   cd [repository-name]
   ```

3. Run the program:

   ```sh
   make run
   ```

   This uses the shortcut defined in the Makefile. Alternatively, you can use:

   ```sh
   go run cmd/unisat-monitor/main.go
   ```

4. (Optional) Create a binary executable:

   ```sh
   make build
   ./bin/projectName
   ```

5. (Optional) Clear the bin directory:

   ```sh
   make clean
   ```

## 🎥 Demo

[asciinema](link)

The bot can catch opportunities with great profit margins: over multiple days of monitoring, the greatest opportunity recorded (in terms of %) was the following, with a listing mistakenly placed at a price **90% lower** than the floor price:
[screenshot](link)

## 🤝 Contributing

Contributions to improve the tool are welcome. Please feel free to submit pull requests or open issues for bugs and feature requests.

## ⚠️ Disclaimer

This tool is for **educational and research purposes only**. Users are responsible for ensuring their use complies with Unisat's terms of service and all applicable laws and regulations.
