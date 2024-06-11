```yaml

Infrastrcuture:
    Timestamp:
        OldestTimestamp: seconds since January 1, 1970 UTC #DONE
        TimeSinceOldest: seconds #DONE

    Controller:Controller[]
        Type:
        Name:
        Location:
            Name:
            Continent:
            Country:
            City:
            GPS_latitude:
            GPS_longitude:
        ServiceLevelAgreement:
        API:
            ICOS_version:
        Any:

    Agent:Cluster[]
        Type: #DONE
        Uuid: #DONE
        Name: #DONE
        ServiceLevelAgreement:
        API:
            ...
        Node[]
            Type: #DONE
            Uuid: #DONE
            Name: #DONE
            Location: #DONE
            Vulnerabilities: #DONE
            ScaScore: #DONE
                ...
            StaticMetrics:
                CPUArchitecture: #DONE
                CPUCores: #DONE
                CPUMaxFrequency: Hertzs #DONE
                GPUCores:
                GPUMaxFrequency:
                GPURAMMemory:
                RAMMemory: Bytes #DONE
                Storage[]
                    Name:
                    Type:
                    Capacity:
            DinamycMetrics:
                Uptime:
                CPUFrequency:
                CPUTemperature: Celsius #DONE
                CPUEnergyConsumption: Joules #DONE
                GPUFrequency:
                GPUTemperature:
                GPUEnergyConsumption:
                FreeRAM: Bytes #DONE
                Storage[]
                    Name:
                    Free:
                Network_usage: # ??
            NetworkInterfaces[]
                Interface_name: #DONE
                Interface_type: #DONE
                Interface_speed: #DONE
                Interface_IP: #DONE
                Interface_status: #DONE
                Interface_subnet_mask: #DONE
                Interface_ingress_usage:
                Interface_egress_usage: 
            Device[]
                Name: #DONE
                Type: #DONE
                Status: #DONE
                API:
                    Communication_protocol:
                    ProtocolVersion:
                    DataFormat:
                    Authentication:
                    Authorization:
        Pod[]
            Name: #DONE
            IP:
            Status: #DONE
            NumberOfContainers: #DONE
            NumberOfApps:
            Container[]
                Name: #DONE
                IP:
                Node: #DONE
                Port:
                ContainerMemory:
                CPUUsage: Percentage #DONE

```
