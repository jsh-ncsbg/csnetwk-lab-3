#include <netinet/in.h>
#include <stdint.h>
#include <sys/socket.h>
#include <unistd.h>

#include <atomic>
#include <cstring>
#include <format>
#include <iostream>
#include <memory>
#include <thread>
#include <vector>

#define SERVER_PORT 8000

typedef struct {
    char name[64];
    int64_t number;
} payload;

int64_t SERVER_NUMBER = 42;
char SERVER_NAME[64] = "Server of the Sokoboys";

class Server {
    int server_socket = -1;
    int server_port;

    std::atomic<bool> is_running = false;
    std::vector<std::thread> workers;

    void handle_client(int client_socket) {
        char buffer[128];
        payload cli_pload;

        ssize_t bytes_received =
            recv(client_socket, buffer, sizeof(buffer) - 1, 0);

        if (!is_running) {
            return;
        }

        if (bytes_received <= 0) {
            if (bytes_received == 0)
                std::cout << "Client disconnected." << std::endl;
            else
                std::cerr << "Failed to receive data from client: "
                          << strerror(errno) << std::endl;
        }

        memcpy(&cli_pload, buffer, sizeof(cli_pload));

        if (cli_pload.number < 1 || cli_pload.number > 100) {
            shutdown();
            close(client_socket);
            return;
        }

        std::cout << std::format("\"{}\" connected to client \"{}\"",
                                 SERVER_NAME, cli_pload.name)
                  << std::endl;
        std::cout << std::format("{} (\"{}\" number) + {} (\"{}\" number) = {}",
                                 SERVER_NUMBER, SERVER_NAME, cli_pload.number,
                                 cli_pload.name,
                                 SERVER_NUMBER + cli_pload.number)
                  << std::endl;

        payload ser_pload;
        memcpy(ser_pload.name, SERVER_NAME, sizeof(payload::name));
        ser_pload.number = SERVER_NUMBER;

        if (send(client_socket, &ser_pload, sizeof(ser_pload), 0) == -1) {
            std::cerr << "Failed to send data to client: " << strerror(errno)
                      << std::endl;
        }

        close(client_socket);
        return;
    }

  public:
    Server(int server_port) : server_port(server_port) {}

    int start() {
        if ((server_socket = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
            std::cerr << strerror(errno) << std::endl;
            return 1;
        }

        sockaddr_in ser_addr = {.sin_family = AF_INET,
                                .sin_port = htons(server_port),
                                .sin_addr = {.s_addr = INADDR_ANY}};

        if (bind(server_socket, (sockaddr *)&ser_addr, sizeof(ser_addr)) ==
            -1) {
            std::cerr << strerror(errno) << std::endl;
            return 1;
        }

        if (listen(server_socket, 5) == -1) {
            std::cerr << strerror(errno) << std::endl;
            return 1;
        }

        std::cout << std::format("Server is listening on port {}...",
                                 server_port)
                  << std::endl;

        is_running = true;

        while (is_running) {
            int client_socket = accept(server_socket, nullptr, nullptr);
            if (!is_running)
                break;

            if (client_socket < 0) {
                std::cerr << strerror(errno) << std::endl;
                continue;
            }

            auto t = std::thread(&Server::handle_client, this, client_socket);
            t.detach();
            workers.push_back(std::move(t));
        }

        for (auto &t : workers) {
            if (t.joinable()) {
                t.join();
            }
        }

        return 0;
    }

    void shutdown() {
        if (server_socket != -1) {
            close(server_socket);
        }
        ::shutdown(server_socket, SHUT_RDWR);
        is_running = false;
    }

    ~Server() {
        std::cout << "Shutting down server..." << std::endl;
        shutdown();
    }
};

int main() {
    Server server(SERVER_PORT);
    server.start();

    server.shutdown();

    return 0;
}