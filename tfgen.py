''' Tensorflow Generator Main '''

import os
import grpc
import argparse

import graphast.nodes as nodes
from graphast.gen import generator, allOps
from tenncorgen.save_ast import save_ast

import proto.serial.data_pb2 as data_pb
import tests.graphmgr.graphmgr_pb2 as graphmgr_pb
import tests.graphmgr.graphmgr_pb2_grpc as graphmgr_rpc

MINDEPTH = os.environ['MINDEPTH'] if 'MINDEPTH' in os.environ else 1 
MAXDEPTH = os.environ['MAXDEPTH'] if 'MAXDEPTH' in os.environ else 10
GRAPH_EXT = ".graph"
REGISTRY = "registry.txt"

def make_graph(outpath):
	rgen = generator("structure.yml", MINDEPTH, MAXDEPTH) 
	root = rgen.generate()
	graphpb = save_ast(root)
	fname = graphpb.gid + GRAPH_EXT
	fpath = os.path.join(outpath, fname)
	print("writing to", fpath)
	with open(fpath, 'wb') as f:
		f.write(graphpb.SerializeToString())
	return graphpb.gid

def main():
	parser = argparse.ArgumentParser(description='Generate Tensorflow script')
	parser.add_argument('--host', dest="host", help='grpc host address')
	parser.add_argument('--out', default='tmp', dest="outpath", help='outpath for generated graphs')
	args = parser.parse_args()
	host = args.host
	outpath = args.outpath

	if not os.path.isdir(outpath):
		os.makedirs(outpath)

	print("generating graph")
	gid = make_graph(outpath)
	with open(os.path.join(outpath, REGISTRY), 'a') as reg:
		reg.write(gid + GRAPH_EXT + '\n')

	print("posting %s to host %s" % (gid, host))
	channel = grpc.insecure_channel(host)
	stub = graphmgr_rpc.GraphmgrStub(channel)
	stub.PostGraph(graphmgr_pb.GraphCreated(gid=gid))

if __name__ == '__main__':
	main()
