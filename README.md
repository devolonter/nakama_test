## ZeptoLab Test Task

Artur Bikmullin

### Environment
-  I initially attempted to launch Nakama with PostgreSQL, but faced some unexpected errors related to an invalid host. As a workaround, I switched to using CockroachDB. I hope this isn't an issue. 
- In order to check the RPC, I added a JavaScript client for Nakama, considering the absence of a Go client library in Nakama currently. The client code is placed in the client directory and necessitates the installation of Node.js and npm.

### Naming Convention and Variable Visibility
- Given that the current plugin is monolithic and all functionalities, except InitModule function, can be private, I consciously chose to keep all elements related to the "public" API accessible. 
- Despite Nakama documentation recommending uppercase for error code constants, I choose Pascal case to adhere to Go's naming convention.

### Solution
- There was some deliberation on whether to utilize composition for In and Out payloads due to their divergent roles.
- Similarly, I was initially uncertain about using string pointers, as it's not a typical Go idiom for simple types and structures. However, as null values were explicitly required in the task, I resorted to pointers.
- The database assignment was perhaps the most challenging to understand and I'm not sure if I executed it correctly. With the task allowing for creative interpretation, I imagine a scenario where we record the function call logs to the database. Although the hard-coded SQL was a compromise I had to make due to time constraints, I would have preferred a different approach.
- cIdeally, the bindPayload should be invoked in the getContent function, but I kept it separate for testing purposes.
- The getContent function returns an error only when the file is not found. This might not be the optimal solution as it can lead to logical errors, for example, when a file is accessible but unreadable.
- Admittedly, testing is not my strongest suit, and my current implementation leaves room for improvement.
- It's disappointing that I couldn't cover tests for all errors. If afforded more time, I would have integrated an abstraction layer to segregate validation, database calls, and other tasks into distinct layers.