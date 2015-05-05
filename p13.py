import string
import hashlib
import socket
import Crypto.Random.random
import base64

"""
The main problem of use the mode CFB is that two strings with the same first
bytes will have the same bytes even encrypted, then I decided to use some
"brute force", well, light brute force, what the script do is the next:
        1. Connect with the server
        2. Sends char by char all the possibilities like: "a", "b", ... until
            have a match, I tested only chars from 31 to 129 ASCII, that are
            the most common
        3. When we have a match, add the char to the decoded string, and start
            agian with the next one :)

Another problem is that the initialization vector is the same, so I could use a
two-pad attack, but this is simplest, fully authomatic for small sentences, and
doesn't requieres to apply a sieve
"""

remote_ip = "54.83.207.93"
port = 12345

def check_proof(proof):
    if len(proof) != 24:
        return False

    h = hashlib.new('sha1')
    h.update(proof)
    if h.digest()[-1:] != '\xff':
        return False

    return True

def auth():
    """ Conect and send the init token generated using brute force """
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    s.connect((remote_ip , port))

    welcome = s.recv(4096)
    base = welcome[28:44]

    while True:
        extra = ''.join(map(lambda _: Crypto.Random.random.choice(string.ascii_letters + string.digits), range(24-len(base))))
        if check_proof(base+extra):
                break

    s.sendall(base+extra)
    s.recv(4096)

    return s

"""
# Lines used for local testing propousals only
s = auth()
text_to_enc1 = "The archaeological and morphological"
s.sendall(text_to_enc1)
encrypted_text = base64.b64decode(s.recv(4096))
s.close()
"""
encrypted_text = base64.b64decode(raw_input())

to_decrypt = ""
while len(to_decrypt) != len(encrypted_text):
    #print "Pos:", len(to_decrypt)+1
    for ascii_char in xrange(31, 129):
        s = auth()
        s.sendall(to_decrypt + chr(ascii_char))
        encrypt_new = base64.b64decode(s.recv(4096))
        s.close()
        if encrypt_new[len(encrypt_new)-1] == encrypted_text[len(encrypt_new)-1]:
                #print "Found:", chr(ascii_char)
                to_decrypt += chr(ascii_char)
                break

print "Result: ", to_decrypt
