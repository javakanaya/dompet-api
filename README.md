# dompet-api
UPDATE:
- Sejauh ini aku udah coba buat bikin routes url /me untuk mengembalikan id, nama, email, password(terenkripsi), dan juga list dompet yang dimiliki. Nah untuk entitasnya aku udah coba buat ada user, dompet, catatan_keuangan, sama detail_user_dompet (hasil tabel many-to-many dari relasi user ke dompet). 
- Oh iya aku juga nambahin kalo setelah user login, itu nanti keluar jwt tokennya di bagian response buat dimasukin ke header nantinya kalau mau akses url secured. (jwt.go)
- Sekalian aku udah coba buat bantu bikin middleware buat authentication, cors + enkripsi password setelah user melakukan register
