import re
import json




data = json.load(open("data/gcp-sku-prices.json"))
regions = json.load(open("data/regions.json"))

machine_types = json.load(open("data/machine-types.json"))

machines = {}

def find_ram_prices(machine_type,region,usageType=None):
    if usageType==None:
        usageType = 'OnDemand'

    family = machine_type.split('-')[0].upper()

    return next(
        (i['pricingInfo'] for i in data['skus'] if i['description'].startswith(family) and i['category']['usageType']==usageType and region in i['serviceRegions']),
        None
    )

def find_cpu_prices(machine_type,region,usageType=None):
    if usageType==None:
        usageType = 'OnDemand'

    family = machine_type.split('-')[0].upper()
    mtype = machine_type.split('-')[1].capitalize()
    resourceGroup = "".join([family,mtype])

    return next(
        (i['pricingInfo'][0] for i in data['skus'] if i['category']['resourceGroup']=='CPU' and i['description'].startswith(family) and i['category']['usageType']==usageType and region in i['serviceRegions']),
        None
    )




for machine in machine_types:
    m = {}
    m['name'] = machine['name']
    m['region'] = "-".join( machine['zone'].split('-')[1:2] )

    cpus_raw = machine['name'].split('-')[-1]
    cpus = int(re.search(r'\d+', cpus_raw)[0])

    import pdb; pdb.set_trace()
    sku_type = ""




    machines[ machine['name'] ] = m
