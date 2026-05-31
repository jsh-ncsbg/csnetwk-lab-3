#include <arpa/inet.h>
#include <netinet/in.h>
#include <stdint.h>
#include <sys/socket.h>
#include <unistd.h>

#include <cstring>
#include <format>
#include <iostream>

typedef struct {
    char name[64];
    int64_t number;
} payload;

std::string get_client_name() {
    while (true) {
        std::cout << "Enter your name (max 63 characters): ";
        std::string name;
        std::getline(std::cin, name);

        if (name.length() > 63) {
            std::cout << "Name too long. Please enter a name with 63 "
                         "characters or fewer."
                      << std::endl;
            continue;
        }

        return name;
    }
}

int64_t get_client_number() {
    while (true) {
        std::cout << "Enter a number: ";
        std::string input;
        std::getline(std::cin, input);

        try {
            int64_t number = std::stoll(input);
            // if (number < 1 || number > 100) {
            //     std::cout << "Please enter a number between 1 and 100." <<
            //     std::endl; continue;
            // }
            return number;
        } catch (const std::invalid_argument &) {
            std::cout << "Invalid input. Please enter a valid number."
                      << std::endl;
        } catch (const std::out_of_range &) {
            std::cout << "Number out of range. Please enter a smaller number."
                      << std::endl;
        }
    }
}

int main() {
    std::string client_name = get_client_name();
    int64_t client_number = get_client_number();

    int cli_sck;
    if ((cli_sck = socket(AF_INET, SOCK_STREAM, 6)) == -1) {
        std::cerr << strerror(errno) << std::endl;
        return 1;
    }

    sockaddr_in ser_addr = {.sin_family = AF_INET, .sin_port = htons(8000)};

    inet_pton(AF_INET, "127.0.0.1", &ser_addr.sin_addr);

    if (connect(cli_sck, (sockaddr *)&ser_addr, sizeof(ser_addr)) == -1) {
        std::cerr << strerror(errno) << std::endl;
        close(cli_sck);
        return 1;
    }

    payload cli_pload;
    memcpy(cli_pload.name, client_name.c_str(), sizeof(payload::name));
    cli_pload.number = client_number;

    if (send(cli_sck, &cli_pload, sizeof(cli_pload), 0) == -1) {
        std::cerr << "Failed to send data to server: " << strerror(errno)
                  << std::endl;
        close(cli_sck);
        return 1;
    }

    payload ser_pload;

    ssize_t bytes_received = recv(cli_sck, &ser_pload, sizeof(ser_pload), 0);

    if (bytes_received < 0) {
        std::cerr << "Failed to receive data from server: " << strerror(errno)
                  << std::endl;
        close(cli_sck);
        return 1;
    }

    std::cout << std::format("{} (\"{}\" number) + {} (\"{}\" number) = {}",
                             ser_pload.number, ser_pload.name, cli_pload.number,
                             cli_pload.name,
                             ser_pload.number + cli_pload.number)
              << std::endl;
}