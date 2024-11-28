# softEng


```mermaid
flowchart TD
    A[Start] --> B[Initialize Blockchain]
    B --> C[Check for Existing Chain]
    C --> |Chain Exists| D[Load Chain]
    C --> |No Chain| E[Create Genesis Block]
    D --> F[Listen for Actions]
    E --> F
    F --> G[Add New Student or Credential]
    G --> |Add Student| H[Store Student in StudentChain]
    G --> |Add Credential| I[Validate Credential]
    I --> |Valid| J[Generate Credential Hash]
    J --> K[Add Credential to Blockchain]
    K --> L[Update Chain]
    L --> M[Broadcast Updated Chain (Optional)]
    M --> F
    M --> N[End]
```