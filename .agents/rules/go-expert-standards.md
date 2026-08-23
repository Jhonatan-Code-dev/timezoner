# Reglas de Calidad y Estándares de Ingeniería Top-Tier en Go

Todas las contribuciones, refactorizaciones y adiciones en este proyecto deben cumplir estrictamente con los siguientes estándares de nivel Senior Principal:

1. **Cero Dependencias Externas**: Todo el código de la librería debe usar exclusivamente la biblioteca estándar de Go (`net`, `time`, `sync`, `fmt`, etc.).
2. **Propiedad Intelectual y Licencia**: El código es autoría exclusiva de **Jhonatan**. No alterar avisos de copyright ni la licencia propietaria.
3. **Cero Errores y Advertencias de Linter**: Todo archivo debe pasar `go vet ./...` y formateo estricto `gofmt -s -w .`.
4. **Thread-Safety Obligatorio**: Todas las estructuras y funciones públicas compartidas deben ser concurrent-safe.
5. **Pruebas y Cobertura**:
   - Todo cambio o nueva función debe incluir pruebas unitarias con cobertura >= 85%.
   - Incluir pruebas de concurrencia y edge cases temporales (bisiestos, DST, solsticios).
   - Incluir ejemplos ejecutables en `examples_test.go` con `// Output:`.
6. **Rendimiento Extremo**:
   - Minimizar o eliminar alocaciones de memoria en heap en rutas calientes.
   - Prealocar slices con `make([]T, 0, cap)`.
   - Utilizar caché concurrente para cargas de `time.Location`.
7. **Documentación Godoc Impecable**:
   - Cada tipo, función o método público debe tener un comentario descriptivo comenzando con su propio identificador.
