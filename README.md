## Security and Best Practices

*   **Security:** Liteforge uses prepared statements to prevent SQL injection. Always use prepared statements when handling user input.
*   **Encryption at Rest (SQLCipher):** If you enable `EncryptAtRest`, ensure you store the `EncryptionKey` securely (e.g., using environment variables or a secrets management solution). Never hardcode the encryption key!
*   **Error Handling:** Handle errors gracefully and provide informative error messages.
*   **Database Migrations:** For production applications, use a database migration tool (e.g., `golang-migrate/migrate`) to manage schema changes.
*   **Input Sanitization:** The `SanitizeInput` function provides minimal sanitization. It's highly recommended to rely on prepared statements for SQL injection prevention and use a dedicated HTML sanitization library (e.g., `github.com/microcosm-cc/bluemonday`) if you need to sanitize HTML content.
*   **Test Thoroughly:** Write comprehensive unit and integration tests to ensure the correctness and reliability of your code.

