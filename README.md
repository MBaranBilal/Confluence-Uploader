# config.json ve client.go Düzenlemesi

"email": "CONFLUENCE_EMAIL" --> bu kısım, Confluence'un cloud mu yoksa server mı olduğuna göre değişebilir. 
- Eğer Cloud ise: `email` kullanılır.
- Eğer Server ise: `username` yazılmalı.
- Alternatif olarak, bu kısım tamamen kaldırılabilir. Bu duruma bağlı olarak client.go içeriğinin de güncellenmesi gerekebilir.