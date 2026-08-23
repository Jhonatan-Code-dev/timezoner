# Autoría, Licencia y Estándares de Documentación Técnica

---

## 1. Autoría Exclusiva y Propiedad Intelectual
- **Autor y Creador Único**: **Jhonatan**.
- **Licencia**: Licencia Propietaria ("All Rights Reserved").
- **Preservación Obligatoria**: Ningún agente ni proceso automatizado puede eliminar, transferir ni alterar la mención de autoría de Jhonatan en los archivos de código fuente, `LICENSE`, ni `README.md`.
- **Prohibición de Copia / Redistribución No Autorizada**: Queda prohibida la redistribución o relicenciamiento no consentido.

## 2. Estándar de Documentación Técnica Limpia (Clean Doc)
- **Cero Clichés de IA**: Prohibido el uso de emojis genéricos de IA (🚀, ✨, 🔥, 💡, ⚡) como decoración superflua en READMEs o documentación técnica.
- **Tono Ejecutivo y Autorizado**: La redacción debe ser formal, precisa, sobria y de nivel senior de ingeniería.
- **Idioma Principal**: Toda la documentación de usuario y `README.md` debe estar escrita en **español** de forma clara, técnica y estructurada.
- **Ejemplos Ejecutables**: Todo ejemplo en la documentación debe compilar y ser directamente verificable mediante código ejecutable en `examples/` o `examples_test.go`.

## 3. Estándar Godoc en Código Fuente
- Cada struct, interfaz, función o constante pública exportada debe contar con su correspondiente bloque de comentarios en español, iniciando siempre con el propio nombre del identificador exportado:
  ```go
  // DBTime es un tipo temporal envoltorio para SQL y JSON que garantiza UTC limpio.
  type DBTime struct { ... }
  ```
