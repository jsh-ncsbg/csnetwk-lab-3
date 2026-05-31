#include <netinet/in.h>
#include <stdint.h>
#include <sys/socket.h>
#include <unistd.h>

#include <atomic>
#include <cstring>
#include <format>
#include <iostream>
#include <thread>

#define SERVER_PORT 8000

typedef struct {
    char name[64];
    int64_t number;
} payload;

int64_t SERVER_NUMBER = 42;
char SERVER_NAME[64] = "Server of the Sokoboys";

int handle_client(const int ser_sck, const int cli_sck) {
    char buffer[128];
    payload cli_pload;

    ssize_t bytes_received = recv(cli_sck, buffer, sizeof(buffer) - 1, 0);

    if (bytes_received <= 0) {
        if (bytes_received == 0)
            std::cout << "Client disconnected." << std::endl;
        else
            std::cerr << "Failed to receive data from client: "
                      << strerror(errno) << std::endl;
    }

    memcpy(&cli_pload, buffer, sizeof(cli_pload));

    // TODO: figure out
    // how to propagate a system exit
    // to all threads and main thread
    // so that graceful shutdown can be achieved
    if (cli_pload.number < 1 || cli_pload.number > 100) {
        close(cli_sck);
        return -1;
    }

    std::cout << std::format("\"{}\" connected to client \"{}\"", SERVER_NAME,
                             cli_pload.name)
              << std::endl;
    std::cout << std::format("{} (\"{}\" number) + {} (\"{}\" number) = {}",
                             SERVER_NUMBER, SERVER_NAME, cli_pload.number,
                             cli_pload.name, SERVER_NUMBER + cli_pload.number)
              << std::endl;

    payload ser_pload;
    memcpy(ser_pload.name, SERVER_NAME, sizeof(payload::name));
    ser_pload.number = SERVER_NUMBER;

    if (send(cli_sck, &ser_pload, sizeof(ser_pload), 0) == -1) {
        std::cerr << "Failed to send data to client: " << strerror(errno)
                  << std::endl;
    }

    close(cli_sck);
    return 0;
}

int main() {
    int ser_sck;

    if ((ser_sck = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
        std::cerr << strerror(errno) << std::endl;
        return 1;
    }

    sockaddr_in ser_addr = {.sin_family = AF_INET,
                            .sin_port = htons(SERVER_PORT),
                            .sin_addr = {.s_addr = INADDR_ANY}};

    if (bind(ser_sck, (sockaddr *)&ser_addr, sizeof(ser_addr)) == -1) {
        std::cerr << strerror(errno) << std::endl;
        close(ser_sck);
        return 1;
    }

    if (listen(ser_sck, 5) == -1) {
        std::cerr << strerror(errno) << std::endl;
        close(ser_sck);
        return 1;
    }

    std::cout << std::format("Server is listening on port {}...", SERVER_PORT)
              << std::endl;

    while (true) {
        int cli_sck = accept(ser_sck, nullptr, nullptr);
        if (cli_sck < 0) {
            std::cerr << strerror(errno) << std::endl;
            continue;
        }
        std::thread(handle_client, ser_sck, cli_sck).detach();
    }

    std::cout << "Shutting down server..." << std::endl;

    close(ser_sck);
}