
#include <iostream>
#include <string>
#include <fstream>
#include <vector>

#include <google/protobuf/descriptor.h>
#include <google/protobuf/descriptor.pb.h>
#include <google/protobuf/dynamic_message.h>
#include <google/protobuf/compiler/importer.h>

uint8_t asc_to_bin(char in)
{
	if (in >= '0' && in <= '9') {
		return in - '0';
	}
	if (in >= 'a' && in <= 'f') {
		return in - 'a' + 10;
	}
	if (in >= 'A' && in <= 'F') {
		return in - 'A' + 10;
	}
	return 0;
}


struct Config
{
	std::string cmd;
	std::string file;
	std::string name;
};

std::map<std::string, Config> g_configs;

std::vector<std::string> split_string(const std::string & input, char sp)
{
	std::vector<std::string> result;

	std::string current;

	for (size_t i = 0; i < input.size(); i++) {
		if (input[i] != sp) {
			current.push_back(input[i]);
		}
		else {
			result.push_back(std::move(current));
			current.clear();
		}
	}

	result.push_back(std::move(current));
	return std::move(result);
}

void load_config_map()
{
	std::ifstream in("./map.conf");
	if (!in) {
		std::cerr << "map.conf not found" << std::endl;
		exit(-10);
	}

	while (in) {
		std::string line;
		getline(in, line);

		std::vector<std::string> params = split_string(line, '\t');
		// std::cerr << params.size() << std::endl;
		if (params.size() != 3) {
			// std::cerr << "map.conf line parse failed:" << line << std::endl;
			continue;
		}

		Config conf;
		conf.cmd = params[0];
		conf.file = params[1];
		conf.name = params[2];

		g_configs[conf.cmd] = conf;
	}

	in.close();
}

int main(int argc, char * argv[])
{
	if (argc != 4) {
		std::cerr << "Usage: " << argv[0] << " command req|rsp bindata" << std::endl;
		return -1;
	}

	std::string command = argv[1];
	command.push_back('_');
	command.append(argv[2]);
	std::string asc_data = argv[3];

	load_config_map();

	Config conf = g_configs[command];
	if (conf.cmd.empty()) {
		std::cerr << "command not found" << std::endl;
		return -1;
	}

	google::protobuf::compiler::DiskSourceTree sourceTree;
	sourceTree.MapPath("", "./pb/");

	google::protobuf::compiler::Importer importer(&sourceTree, NULL);
	importer.Import(conf.file);

	const google::protobuf::Descriptor * descriptor = importer.pool()->FindMessageTypeByName(conf.name);

	google::protobuf::DynamicMessageFactory factory;
	const google::protobuf::Message * message = factory.GetPrototype(descriptor);

	google::protobuf::Message * instance = message->New();

	std::string data;

	for (int i = 0; i < asc_data.size(); i+=2) {
		uint8_t c = (asc_to_bin(asc_data[i]) << 4) | asc_to_bin(asc_data[i+1]);
		data.push_back((char)c);
	}

	if (!instance->ParseFromString(data)) {
		std::cerr << "ParseFromArray failed" << std::endl;
		return -2;
	}

	std::cout << "Succ:" << std::endl;
	std::cout << instance->DebugString() << std::endl;
	return 0;
}

