''' Tensorflow Generator Main '''

import sys
import os
import grpc
import argparse

from graphast.gen import generator, allOps
from tenncorgen.save_ast import save_ast
from tenncorgen.tfgen import tf_gen

import proto.serial.data_pb2 as data_pb
import tests.graphmgr.graphmgr_pb2 as graphmgr_pb
import tests.graphmgr.graphmgr_pb2_grpc as graphmgr_rpc

MINDEPTH = os.environ['MINDEPTH'] if 'MINDEPTH' in os.environ else 1 
MAXDEPTH = os.environ['MAXDEPTH'] if 'MAXDEPTH' in os.environ else 10
PY_EXT = ".py"
GRAPH_EXT = ".graph"

def make_graph(outpath, root, graphpb):
	script = tf_gen(root, graphpb.gid, graphpb.create_order,
		out_prefix=outpath + "/",
		external='external/com_github_mingkaic_tenncor')

	gpath = os.path.join(outpath, graphpb.gid + GRAPH_EXT)
	with open(gpath, 'wb') as f:
		f.write(graphpb.SerializeToString())

	spath = os.path.join(outpath, graphpb.gid + PY_EXT)
	with open(spath, 'w') as f:
		f.write(script)

def graph_created(gid, channel):
	stub = graphmgr_rpc.GraphmgrStub(channel)
	stub.PostGraph(graphmgr_pb.GraphCreated(gid=gid))

def main():
	print("running tfgen")
	parser = argparse.ArgumentParser(description='Generate Tensorflow script')
	parser.add_argument('--host', dest="host", help='grpc host address')
	parser.add_argument('--out', dest="outpath", help='outpath for generated graphs')
	parser.add_argument('--rando', action='store_true', help='generate random graphs')
	args = parser.parse_args()
	host = args.host
	outpath = args.outpath
	assert(host is not None)
	assert(outpath is not None)

	if not os.path.isdir(outpath):
		os.makedirs(outpath)

	print("posting graphs generated to host %s" % host)
	channel = grpc.insecure_channel(host)
	print("writing graph and script to", outpath)

	if args.rando:
		print("generating random graph")
		rgen = generator("structure.yml", MINDEPTH, MAXDEPTH) 
		root = rgen.generate()
		graphpb = save_ast(root)
		make_graph(outpath, root, graphpb)
		gid = graphpb.gid
		graph_created(gid, channel)
	else:
		# todo: check for all ops existence
		print("generating graph for each op")
		ops = allOps("structure.yml")
		for opname in ops:
			tops = ops[opname]
			for i in range(len(tops)):
				root = tops[i]
				newGid = opname
				if i > 0:
					newGid = newGid + str(i)
				graphpb = save_ast(root)
				graphpb.gid = newGid
				make_graph(outpath, root, graphpb)
				graph_created(newGid, channel)

if __name__ == '__main__':
	main()
