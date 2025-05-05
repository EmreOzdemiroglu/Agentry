# Agentry 🤖

Agentry is a command-line interface (CLI) tool designed for creating, managing, and executing chains of local AI agents powered by Ollama. Build complex workflows by linking different agents together, experiment with prompts, and observe the agent interactions step-by-step.

## Current Features ✨

*   **Interactive Chain Creation:** Define sequences of agents, specifying unique names, Ollama models, system prompts, and user prompt templates for each step.
*   **Agent Reusability:** Reuse agent definitions (model, system prompt) within a chain, optionally modifying the user prompt for specific steps.
*   **Chain Management:**
    *   Save agent chain configurations to JSON files (stored in the `Agentries/` directory).
    *   Load existing chains for running or editing.
    *   Edit chains by adding or deleting agents.
*   **Execution Modes:**
    *   **Sequential:** Run agents one after another, passing the output of one as the input to the next.
    *   **Looping:** Optionally configure chains to loop, feeding the final output back to the first agent for iterative processing.
*   **Detailed Logging:** Capture the entire execution flow, including inputs, prompts, responses, and timings for each agent, into timestamped Markdown files (stored in the `chatruns/` directory).
*   **Local First:** Operates entirely with your local Ollama instance and models.

## Getting Started 🚀

1.  **Prerequisites:**
    *   Go (version 1.21 or later recommended).
    *   Ollama installed and running with desired models downloaded (e.g., `ollama pull gemma:2b`).
2.  **Clone the repository (if applicable):**
    ```bash
    git clone <your-repo-url>
    cd agentry
    ```
3.  **Run the application:**
    ```bash
    go run .
    ```
4.  **Follow the CLI prompts:**
    *   `[C]`reate a new chain.
    *   `[L]`oad an existing chain from the `Agentries/` folder.
    *   `[Q]`uit.

## The Vision: The Ultimate Agentic Toolkit 🌌

Agentry aims to evolve into a comprehensive and flexible platform for building and orchestrating sophisticated AI agent systems. The current CLI is just the beginning!

**Upcoming Focus & Future Roadmap:**

*   **Enhanced CLI:** Continue refining the core CLI experience.
*   **Parallel Execution:** Enable agents within a chain (or across chains) to run concurrently for faster processing and more complex workflows.
*   **Agent Tooling:** Allow agents to interact with external tools, APIs, and data sources.
*   **Search Integration:** Equip agents with capabilities to perform web searches, document lookups, etc.
*   **Advanced Chain Orchestration:**
    *   Combine multiple saved chains.
    *   Create branching logic and conditional execution paths.
    *   Dynamically route information between agents and chains.
*   **Visualization:** Generate visual representations of agent chains (e.g., using JSON to Mermaid diagrams viewable in tools like Excalidraw) to better understand their structure and flow.
*   **Obsidian Plugin:** Integrate Agentry directly into Obsidian for seamless knowledge management and agent-driven content creation/analysis within your notes.

The goal is to make Agentry the most powerful, flexible, and user-friendly tool for anyone working with local LLMs and agentic systems.

## Contributing 🤝

Contributions are welcome! (Details on contribution guidelines will be added soon).

---

*Build the future of local AI agents with Agentry!* 