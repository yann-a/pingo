nbgroups = 0

"""
Transforme un fichier de règles en une regex pas jolie jolie.
"""
def compile(rule):
    global nbgroups

    # Syntactic sugar for sequential rules
    if isinstance(rule, dict):
        ret = ""
        for key in rule:
            ret += compile([key, rule[key]])

        return ret

    if isinstance(rule, int):
        return str(rule)+"\\n"

    if rule[0] == "select":
        regexps = map(compile, rule[1])

        return "(?:" + "|".join(regexps) + ")"

    if rule[0] == "sequential":
        return "".join(map(compile, rule[1]))

    # rule[0] == "parallel"
    # /!\ Cette règle ne supporte que des sous règles étant directement des entiers distincts

    regexps = list(map(compile, rule[1]))
    inner = "(?:" + "|".join(regexps) + ")"
    ret = ""
    for i in range(len(regexps)-1):
        nbgroups += 1
        gpName = "g"+str(nbgroups)
        ret += "(?:(?P<"+gpName+">"+inner+")(?!(?:[^s]\\n){0,"+str(len(regexps)-i-2)+"}(?P="+gpName+")))"

    ret += inner

    return ret

import glob
import re
import yaml
import subprocess

for file in glob.glob('tests/pi/*.pi'):
    result = subprocess.run(['./pingo', file], stdout=subprocess.PIPE, stderr=subprocess.DEVNULL)

    with open(file[:-3] + ".expect", 'r') as stream:
        try:
            rule = yaml.safe_load(stream)
            nbgroups = 0
            regex = "^"+compile(rule)+"$"

            if re.match(regex, result.stdout.decode("utf-8")):
                print("ok")
            else:
                print("error while executing "+file)
        except yaml.YAMLError as exc:
            print(exc)


for file in glob.glob('tests/lambda/*.lambda'):
    result = subprocess.run(['./pingo', '-translate', file], stdout=subprocess.PIPE, stderr=subprocess.DEVNULL)

    with open(file[:-7] + ".expect", 'r') as stream:
        if stream.read() == result.stdout.decode("utf-8"):
            print("ok")
        else:
            print("error while executing "+file)
