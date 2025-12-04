import webview
import socket
import time
import subprocess
import sys
import requests

host = "localhost"
port = 8080

if __name__ == '__main__':
    proccess = subprocess.Popen(["go", "run", "."])

    while True:
        try:
            response = requests.get(f'http://{host}:{port}/check')
            if response.status_code == 200:
                break
        except:
            time.sleep(0.5)
    
    webview.settings['ALLOW_DOWNLOADS'] = True
    webview.create_window("Douglas", f'http://{host}:{port}/Home')
    webview.start()

    if sys.platform == "win32":
        proccess.terminate()
    else:
        proccess.kill()

    proccess.wait()