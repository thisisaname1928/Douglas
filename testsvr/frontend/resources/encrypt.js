let privateKey
let publicKey

async function generateECDHKeys() {
    const ec = new CryptoJS.elliptic.ec('p256');

    const key = ec.genKeyPair();
    publickey = key.getPublic().encodeCompressed('hex');
}

async function initEncrypting() {
    generateECDHKeys()
}

initEncrypting()