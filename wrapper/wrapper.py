import webview
import socket
import time
import subprocess

host = "localhost"
port = 8080

if __name__ == '__main__':
    proccess = subprocess.Popen(["go", "run", "."])
    
    webview.settings['ALLOW_DOWNLOADS'] = True
    webview.create_window("Douglas", f'http://localhost:8080/Home')
    webview.start()

    proccess.terminate()
    proccess.kill()
    time.sleep(0.5)