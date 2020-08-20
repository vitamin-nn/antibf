package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/vitamin-nn/otus_anti_bruteforce/internal/config"
	grpcService "github.com/vitamin-nn/otus_anti_bruteforce/internal/grpc"
	"google.golang.org/grpc"
)

type modifyFunc func(context.Context, *grpcService.ModifyListRequest, ...grpc.CallOption) (*grpcService.ModifyResponse, error)

func cliCmd(cfg *config.Config) *cobra.Command {
	cliCmd := &cobra.Command{
		Use:   "cli",
		Short: "Command line interface for managing of antibf service",
		Run: func(cmd *cobra.Command, args []string) {
			log.Fatalf("Please specify cli-command")
		},
	}
	clearCmd := getClearCmd(cfg)
	cliCmd.AddCommand(clearCmd)

	addWhiteCmd := getAddWhiteCmd(cfg)
	cliCmd.AddCommand(addWhiteCmd)

	removeWhiteCmd := getRemoveWhiteCmd(cfg)
	cliCmd.AddCommand(removeWhiteCmd)

	addBlackCmd := getAddBlackCmd(cfg)
	cliCmd.AddCommand(addBlackCmd)

	removeBlackCmd := getRemoveBlackCmd(cfg)
	cliCmd.AddCommand(removeBlackCmd)

	return cliCmd
}

func getClearCmd(cfg *config.Config) *cobra.Command {
	var ip string
	var login string

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clears specified bucket",
		Run: func(cmd *cobra.Command, args []string) {
			grpcConn := getGrpcConn(cfg.GrpcServer.Addr)
			defer grpcConn.Close()
			grpcClient := grpcService.NewAntiBruteforceServiceClient(grpcConn)

			req := &grpcService.ClearRequest{}
			if login != "" {
				req.Login = login
			}
			if ip != "" {
				req.Ip = ip
			}

			if ip == "" && login == "" {
				log.Fatalln("Must be set at least one param")
			}

			resp, err := grpcClient.Clear(context.Background(), req)
			if err != nil {
				log.Fatalf("Grpc request transport error: %v", err)
			}
			processGrpcResponse(resp)
		},
	}
	cmd.Flags().StringVar(&login, "login", "", "Login that will cleared")
	cmd.Flags().StringVar(&ip, "ip", "", "IP that will cleared")

	return cmd
}

func getAddWhiteCmd(cfg *config.Config) *cobra.Command {
	var ip string
	cmd := &cobra.Command{
		Use:   "addwhite",
		Short: "Adds ip network to the white list",
		Run: func(cmd *cobra.Command, args []string) {
			grpcConn := getGrpcConn(cfg.GrpcServer.Addr)
			defer grpcConn.Close()
			grpcClient := grpcService.NewAntiBruteforceServiceClient(grpcConn)
			modifyList(grpcClient.AddToWhiteList, ip)
		},
	}
	cmd.Flags().StringVar(&ip, "ip", "", "IP network")
	err := cmd.MarkFlagRequired("ip")
	if err != nil {
		log.Fatalf("Marking flag required error: %v", err)
	}

	return cmd
}

func getRemoveWhiteCmd(cfg *config.Config) *cobra.Command {
	var ip string
	cmd := &cobra.Command{
		Use:   "rmwhite",
		Short: "Remove ip network from the white list",
		Run: func(cmd *cobra.Command, args []string) {
			grpcConn := getGrpcConn(cfg.GrpcServer.Addr)
			defer grpcConn.Close()
			grpcClient := grpcService.NewAntiBruteforceServiceClient(grpcConn)
			modifyList(grpcClient.RemoveFromWhiteList, ip)
		},
	}
	cmd.Flags().StringVar(&ip, "ip", "", "IP network")
	err := cmd.MarkFlagRequired("ip")
	if err != nil {
		log.Fatalf("Marking flag required error: %v", err)
	}

	return cmd
}

func getAddBlackCmd(cfg *config.Config) *cobra.Command {
	var ip string
	cmd := &cobra.Command{
		Use:   "addblack",
		Short: "Adds ip network to the black list",
		Run: func(cmd *cobra.Command, args []string) {
			grpcConn := getGrpcConn(cfg.GrpcServer.Addr)
			defer grpcConn.Close()
			grpcClient := grpcService.NewAntiBruteforceServiceClient(grpcConn)
			modifyList(grpcClient.AddToBlackList, ip)
		},
	}
	cmd.Flags().StringVar(&ip, "ip", "", "IP network")
	err := cmd.MarkFlagRequired("ip")
	if err != nil {
		log.Fatalf("Marking flag required error: %v", err)
	}

	return cmd
}

func getRemoveBlackCmd(cfg *config.Config) *cobra.Command {
	var ip string
	cmd := &cobra.Command{
		Use:   "rmblack",
		Short: "Remove ip network from the black list",
		Run: func(cmd *cobra.Command, args []string) {
			grpcConn := getGrpcConn(cfg.GrpcServer.Addr)
			defer grpcConn.Close()
			grpcClient := grpcService.NewAntiBruteforceServiceClient(grpcConn)
			modifyList(grpcClient.RemoveFromBlackList, ip)
		},
	}
	cmd.Flags().StringVar(&ip, "ip", "", "IP network")
	err := cmd.MarkFlagRequired("ip")
	if err != nil {
		log.Fatalf("Marking flag required error: %v", err)
	}

	return cmd
}

func getGrpcConn(grpcAddr string) *grpc.ClientConn {
	grpcConn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	return grpcConn
}

func processGrpcResponse(resp *grpcService.ModifyResponse) {
	res, ok := resp.Result.(*grpcService.ModifyResponse_Success)
	if !ok {
		res, ok := resp.Result.(*grpcService.ModifyResponse_Error)
		if !ok {
			log.Fatal("Unknown grpc response")
		}
		log.Fatalf("Grpc request error: %s", res.Error)
	}
	fmt.Printf("Success: %v\n", res.Success)
}

func modifyList(f modifyFunc, ip string) {
	req := &grpcService.ModifyListRequest{
		Ip: ip,
	}
	resp, err := f(context.Background(), req)
	if err != nil {
		log.Fatalf("Grpc request transport error: %v", err)
	}
	processGrpcResponse(resp)
}
