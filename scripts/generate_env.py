import configparser

def parse_config_to_env(config_file, env_file):
    config = configparser.ConfigParser()
    config.read(config_file)

    with open(env_file, 'w') as env:
        for section in config.sections():
            # Replace section headers to environment-friendly format
            section_name = section.replace('.', '_').upper()

            # Write section header comment
            # env.write(f"# === {section_name} ===\n")

            for key, value in config.items(section):
                # Convert key to uppercase and join it with section name
                env_var_name = f"{section_name}_{key}".upper()

                try:
                    tmp = value.split("#")
                    value = tmp[0].strip()
                except Exception:
                    pass

                # Write variable in the form VAR=VALUE
                env.write(f"{env_var_name}={value}\n")
            
            env.write("\n")

if __name__ == "__main__":
    parse_config_to_env("config.toml", ".env")
