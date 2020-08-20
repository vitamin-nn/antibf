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

//var grpcConn *grpc.ClientConn
//var grpcClient grpcService.AntiBruteforceServiceClient

type modifyFunc func(context.Context, *grpcService.ModifyListRequest, ...grpc.CallOption) (*grpcService.ModifyResponse, error)

var ip string

func cliCmd(cfg *config.Config) *cobra.Command {
	var login string

	cliCmd := &cobra.Command{
		Use:   "cli",
		Short: "Command line interface for managing of antibf service",
		Run: func(cmd *cobra.Command, args []string) {
			log.Fatalf("Please specify cli-command")
		},
	}

	clearCmd := &cobra.Command{
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
			resp, err := grpcClient.Clear(context.Background(), req)
			if err != nil {
				log.Fatalf("Grpc request transport error: %v", err)
			}
			processGrpcResponse(resp)
		},
	}
	clearCmd.Flags().StringVar(&login, "login", "", "Login that will cleared")
	clearCmd.Flags().StringVar(&ip, "ip", "", "IP that will cleared")
	cliCmd.AddCommand(clearCmd)

	addWhiteCmd := &cobra.Command{
		Use:   "addwhite",
		Short: "Adds ip network to the white list",
		Run: func(cmd *cobra.Command, args []string) {
			grpcConn := getGrpcConn(cfg.GrpcServer.Addr)
			defer grpcConn.Close()
			grpcClient := grpcService.NewAntiBruteforceServiceClient(grpcConn)
			modifyList(grpcClient.AddToWhiteList, ip)
		},
	}
	addIpFlag(addWhiteCmd)
	cliCmd.AddCommand(addWhiteCmd)

	removeWhiteCmd := &cobra.Command{
		Use:   "rmwhite",
		Short: "Remove ip network from the white list",
		Run: func(cmd *cobra.Command, args []string) {
			grpcConn := getGrpcConn(cfg.GrpcServer.Addr)
			defer grpcConn.Close()
			grpcClient := grpcService.NewAntiBruteforceServiceClient(grpcConn)
			modifyList(grpcClient.RemoveFromWhiteList, ip)
		},
	}
	addIpFlag(removeWhiteCmd)
	cliCmd.AddCommand(removeWhiteCmd)

	addBlackCmd := &cobra.Command{
		Use:   "addblack",
		Short: "Adds ip network to the black list",
		Run: func(cmd *cobra.Command, args []string) {
			grpcConn := getGrpcConn(cfg.GrpcServer.Addr)
			defer grpcConn.Close()
			grpcClient := grpcService.NewAntiBruteforceServiceClient(grpcConn)
			modifyList(grpcClient.AddToBlackList, ip)
		},
	}
	addIpFlag(addBlackCmd)
	cliCmd.AddCommand(addBlackCmd)

	removeBlackCmd := &cobra.Command{
		Use:   "rmblack",
		Short: "Remove ip network from the black list",
		Run: func(cmd *cobra.Command, args []string) {
			grpcConn := getGrpcConn(cfg.GrpcServer.Addr)
			defer grpcConn.Close()
			grpcClient := grpcService.NewAntiBruteforceServiceClient(grpcConn)
			modifyList(grpcClient.RemoveFromBlackList, ip)
		},
	}
	addIpFlag(removeBlackCmd)
	cliCmd.AddCommand(removeBlackCmd)

	return cliCmd
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

func addIpFlag(cmd *cobra.Command) {
	cmd.Flags().StringVar(&ip, "ip", "", "IP network")
	err := cmd.MarkFlagRequired("ip")
	if err != nil {
		log.Fatalf("Marking flag ruquired error: %v", err)
	}
}
