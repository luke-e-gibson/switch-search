{
  description = "Switch-Search";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      # nix build
      packages.${system}.default = pkgs.buildGoModule {
        pname = "switch-search";
        version = "0.1.0";
        src = ./.;
        vendorHash = "sha256-dyRlT2EfC7iJqwFkvmmt9VTCYc2uBWedKZGo7yqMOuM="; # replace after first build, see below
      };  

      # nix run
      apps.${system}.default = {
        type = "app";
        program = "${self.packages.${system}.default}/bin/switch-search";
      };

      # systemd service
      nixosModules.default = { config, lib, ... }:
        let cfg = config.services.switch-search;
        in {
          options.services.switch-search = {
            enable = lib.mkEnableOption "Switch-Search server";
            openFirewall = lib.mkEnableOption "opening port 9080 in the firewall";
          };

          config = lib.mkIf cfg.enable {
            systemd.services.switch-search = {
              description = "Switch-Search";
              wantedBy = [ "multi-user.target" ];
              after = [ "network.target" ];
              serviceConfig = {
                ExecStart = "${self.packages.${system}.default}/bin/switch-search";
                DynamicUser = true;
                Restart = "on-failure";
              };
            };

          };
        };

      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [ go gopls gotools ];
      };
    };
}