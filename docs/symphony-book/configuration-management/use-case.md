# Configuration use cases

Symphony configuration supports different kind of expressions. Here are some samples.

## Referring parent properties
```
model.CatalogState{
    Id: "catalog1",
    Spec: &model.CatalogSpec{
        ParentName: "parent",
        Properties: map[string]interface{}{
            "components": []model.ComponentSpec{
                {
                    Name: "name",
                    Type: "type",
                },
            },
        },
    },
}

model.CatalogState{
    Id: "parent",
    Spec: &model.CatalogSpec{
        Properties: map[string]interface{}{
            "parentAttribute": "This is father",
        },
    },
}

Read("catalog1", "parentAttribute")
"This is father"
```

## Referring deployment spec params
```
utils.EvaluationContext{
    DeploymentSpec: model.DeploymentSpec{
        Instance: model.InstanceSpec{
            Solution: "fake-solution",
            Arguments: map[string]map[string]string{
                "component-1": {
                    "a": "new-value",
                },
            },
        },
        SolutionName: "fake-solution",
        Solution: model.SolutionSpec{
            Components: []model.ComponentSpec{
                {
                    Name: "component-1",
                    Parameters: map[string]string{
                        "a": "b",
                        "c": "d",
                    },
                },
            },
        },
    },
    Component: "component-1",
}
${{$param(a)}}
"new-value"
```

## Referring evaluation context property
```
utils.EvaluationContext{
    Properties: map[string]string{
        "foo": "bar",
    },
}
${{$property(foo)}}
"bar"
```

## Referring evaluation context input
```
utils.EvaluationContext{
    Inputs: map[string]interface{}{
        "foo":  "bar",
        "book": "title",
    },
}
"${{$input(foo)}}"
"bar"
```

## Referring evaluation context output
```
utils.EvaluationContext{
    Outputs: map[string]map[string]interface{}{
        "foo": map[string]interface{}{
            "bar": 5,
        },
    },
}
"${{$output(foo,bar)}}"
5
```

## Referring evaluation context value and context
```
utils.EvaluationContext{
    Value: 6,
}
${{$and($gt($val(),5), $lt($val(),10))}}"
True

utils.EvaluationContext{
    Value: map[string]interface{}{
        "foo": map[string]interface{}{
            "bar": "baz",
        },
    },
}
"${{$equal($val('$.foo.bar'), "baz")}}"
True

"${{$equal($context('$.foo.bar'), "baz")}}"
True
```

## Referring evaluation context deployment instance
```
utils.EvaluationContext{
    DeploymentSpec: model.DeploymentSpec{
        Instance: model.InstanceSpec{
            Name: "instance-1",
        },
        SolutionName: "fake-solution",
        Solution: model.SolutionSpec{
            Components: []model.ComponentSpec{
                {
                    Name: "component-1",
                    Parameters: map[string]string{
                        "a": "b",
                        "c": "d",
                    },
                },
            },
        },
    },
}
"${{$if($eq($instance(), "instance-1"), "do some thing", "do nothing")}}"
"do some thing"
```

## Referring $config to evaluate other catalogs
```
model.CatalogState{
    Id: "a",
    Spec: &model.CatalogSpec{
        ParentName: "a",
        Properties: map[string]interface{}{
            "b": 4,
            "c": false
        },
    },
}

utils.EvaluationContext{
    ConfigProvider: configProvider,
    DeploymentSpec: model.DeploymentSpec{
        Instance: model.InstanceSpec{
            Solution: "fake-solution",
            Arguments: map[string]map[string]string{
                "component-1": {
                    "a": "new-value",
                },
            },
        },
        SolutionName: "fake-solution",
        Solution: model.SolutionSpec{
            Components: []model.ComponentSpec{
                {
                    Name: "component-1",
                    Properties: map[string]interface{}{
                        "foo": "${{$if($between($config(a,b), 1, 5), "footrue", "foofalse")}}",
                        "bar": "${{$if($not($config(a,c)), "barfalse", "bartrue")}}",
                    },
                },
            },
        },
    },
    Component: "component-1",
}
DeploymentSpec.Solution.Components[0].Properties["foo"]
"footrue"

DeploymentSpec.Solution.Components[0].Properties["bar"]
"barfalse"
```