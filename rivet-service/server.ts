import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { runGraph, loadProjectFromFile, type Project } from '@ironclad/rivet-node';
import * as path from 'path';
import * as fs from 'fs';
import { fileURLToPath } from 'url';
import { dirname } from 'path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const PROTO_PATH = path.join(__dirname, 'rivet.proto');

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

const rivetProto = grpc.loadPackageDefinition(packageDefinition).rivet as any;

async function executeGraph(call: any, callback: any) {
  const { graph_id, inputs, toolbelt_json, project_content } = call.request;

  console.log(`[Rivet Service] Executing graph: ${graph_id}`);

  try {
    let project;
    
    if (project_content) {
      // Load project from request content
      project = JSON.parse(project_content) as Project;
    } else {
      // Fallback to a default project file if it exists
      const defaultProjectPath = path.join(__dirname, 'project.rivet-project');
      if (fs.existsSync(defaultProjectPath)) {
        project = await loadProjectFromFile(defaultProjectPath);
      } else {
        throw new Error("No project content provided and no default project found");
      }
    }

    // Convert Struct inputs to Rivet inputs
    // Note: Rivet expects inputs as Record<string, DataValue>
    // We need to map simple JSON values to Rivet DataValues
    const rivetInputs: Record<string, any> = {};
    if (inputs && inputs.fields) {
        for (const [key, val] of Object.entries(inputs.fields)) {
            rivetInputs[key] = {
                type: 'string', // Simplified assumption for now
                value: (val as any).stringValue || (val as any).numberValue || JSON.stringify(val)
            };
        }
    }

    // Run the graph
    const result = await runGraph(project, {
      graph: graph_id,
      inputs: rivetInputs,
      // remoteDebugger: true, // Optional: Enable for debugging with Rivet UI
    });

    // Convert outputs back to Struct
    const outputStruct: Record<string, any> = {};
    for (const [key, val] of Object.entries(result)) {
        outputStruct[key] = val.value;
    }

    callback(null, { outputs: outputStruct });
  } catch (error: any) {
    console.error('[Rivet Service] Error:', error);
    callback({
      code: grpc.status.INTERNAL,
      details: error.message,
    }, null);
  }
}

function main() {
  const server = new grpc.Server();
  server.addService(rivetProto.RivetService.service, { ExecuteGraph: executeGraph });
  
  const address = '0.0.0.0:50051';
  server.bindAsync(address, grpc.ServerCredentials.createInsecure(), (err, port) => {
    if (err) {
      console.error(`[Rivet Service] Failed to bind: ${err}`);
      return;
    }
    console.log(`[Rivet Service] listening on ${address} (port ${port})`);
    // server.start(); // start() is deprecated in newer grpc-js and not needed if bindAsync is successful? 
    // Actually, checking docs: server.start() IS deprecated but might still be needed in some versions?
    // The error "server must be bound in order to start" happens if we call start() before binding is complete.
    // But here we are IN the callback.
    // However, if the error happens, we shouldn't call start.
    // If no error, it is bound.
  });
}

main();
