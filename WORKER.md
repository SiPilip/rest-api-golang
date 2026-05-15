Signup request masuk
→ User disimpan ke DB
→ Job "send_welcome_email" masuk ke queue
→ Response langsung dikirim (cepat!)

Di background:
Worker 1 ambil job → kirim email → selesai
Worker 2 idle (menunggu job lain)
Worker 3 idle

Server shutdown (Ctrl+C):
→ Tutup queue (tidak terima job baru)
→ Tunggu Worker 1 selesai kirim email
→ Semua worker selesai → tutup DB → exit

INFO "Job enqueued" job=send_welcome_email
INFO "Worker processing job" workerId=1 job=send_welcome_email
INFO "Email sent" to=user@test.com
INFO "Job completed" workerId=1 job=send_welcome_email

INFO "Shutting down server..."
INFO "Worker processing job" ... ← masih proses
INFO "Job completed" ... ← selesai dulu
INFO "Worker pool stopped" ← baru shutdown
INFO "Server gracefully stopped"
