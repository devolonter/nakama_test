import { Client } from '@heroiclabs/nakama-js'
import XMLHttpRequest from 'xhr2'

global.XMLHttpRequest = XMLHttpRequest

const client = new Client('defaultkey', '127.0.0.1', 7350, false)
const session = await client.authenticateDevice('my-test-device-id', true)

console.log(await client.rpc(session, 'get_content', {}))