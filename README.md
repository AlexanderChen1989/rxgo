# rxgo
Reactive Extensions implementation for the Go Programming Language.

## Goal
* [ ] 实现ReactiveX规范
* [ ] 符合Go的特色
* [ ] 可控的资源消耗
* [ ] 高效

## Design
```
value -> observable -> observable -> observer
channel + goroutine

goroutine -> chan -> (do compute) -> chan -> observer
```

## TODO
