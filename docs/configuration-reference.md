| Variable           | Required | Description                                                                                       | Example                               |
|--------------------|----------|---------------------------------------------------------------------------------------------------|---------------------------------------|
| PVE_URL            | Yes      | The base Proxmox VE API URL structure - https://<your-server-ip-or-domain>:8006/api2/json/        | https://192.168.18.999:8006/api2/json |
| PVE_USER           | Yes      | The name of the Proxmox user responsible for the API communication                                | Orcha                                 |
| PVE_REALM          | Yes      | The user's realm (PVE)                                                                            | PVE                                   |
| PVE_USER_REALM     | Yes      | User + Realm -> PVE_USER@PVE_REALM                                                                | Orcha@PVE                             |
| PVE_TOKEN_ID       | Yes      | Proxmox VE API Token ID                                                                           | api-token                             |
| PVE_TOKEN          | Yes      | Proxmox VE API Token                                                                              | aaa99aaa-a9a9-99aa-aaa9-a99999a999a9  |
| PVE_NODE_NAME      | Yes      | The name of the node which will be accessed by Orcha                                              | main                                  |
| SSH_KEY_FILE       | Yes      | Path to the SSH Key, in the machine running Orcha, that will be used to access the Proxmox Server | "/home/apiuser/.ssh/apiKey"           |
| SSH_KEY_PASSPHRASE | Yes      | The SSH Key passphrase                                                                            | 999aaa                                |
| SSH_PVE_IP         | Yes      | The main Proxmox IP address                                                                       | 192.168.18.999                        |
| SSH_PVE_PORT       | Yes      | The main Proxmox SSH port                                                                         | 22                                    |