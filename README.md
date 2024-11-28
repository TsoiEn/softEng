# softEng


```mermaid
flowchart TD
    A[Start] --> B[Initialize Blockchain]
    B --> C[Check for Existing Chain]
    C --> |Chain Exists| D[Load Chain]
    C --> |No Chain| E[Create Genesis Block]
    D --> F[Validate Chain]
    E --> F
    F --> G[Listen for Transactions / add credential]
    G --> |Add Credential| H[Validate Credential]
    H --> |Valid| I[Generate Credential Hash]
    I --> J[Add Credential to Blockchain]
    J --> K[Update Chain]
    K --> L[Broadcast Updata chain]
    L --> G
    L --> M[End]
```